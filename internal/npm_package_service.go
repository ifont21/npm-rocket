package internal

import (
	"fmt"

	"github.com/ifont21/pre-releaser-cli/internal/git_util"
)

type NPMPackageService struct {
	pkgJsonRepository   PkgJSONRepository
	changelogRepository ChangeLogRepository
	gitRepository       GitRepository
	ghRepository        GithubRepository
}

func NewNPMPackageService(pkgJsonRepository PkgJSONRepository, changelogRepository ChangeLogRepository, gitRepository GitRepository) *NPMPackageService {
	return &NPMPackageService{
		pkgJsonRepository:   pkgJsonRepository,
		changelogRepository: changelogRepository,
		gitRepository:       gitRepository,
		ghRepository:        NewGithubRepositoryImpl(),
	}
}

func (n *NPMPackageService) BumpNPMPackage(filePath string, bumpType string) (string, error) {
	version, err := n.pkgJsonRepository.BumpNPMPackage(filePath, bumpType)
	if err != nil {
		return "", err
	}

	return version, nil
}

// function to bump and update the CHANGELOG.md file
func (n *NPMPackageService) BumpNPMPackagesAndChangelog(repoPath string, libs []string, commits string) error {
	listAffected, err := n.gitRepository.GetAffectedLibsFromGitCommits(repoPath, "")
	if err != nil {
		fmt.Println("Error getting listAffected", err)
		return err
	}
	fmt.Println("listAffected -->", listAffected)
	commitMessages, err := git_util.GetCommitMessagesByDate("", repoPath)
	if err != nil {
		fmt.Println("Error getting commitMessages", err)
		return err
	}

	// loop through the libs and bump the versions
	for _, lib := range listAffected {
		fmt.Println("Bumping package... ", lib)
		libPath := fmt.Sprintf("%s/libs/%s", repoPath, lib)
		fmt.Println("commit messages: ", commitMessages)
		// determine bump type
		commitsByScope, err := n.gitRepository.SummarizeCommitsByScope(lib, commitMessages)
		if err != nil {
			fmt.Println("Error getting commits by scope", err)
			return err
		}
		fmt.Printf("text selected by the scope `%s`\n%s", lib, commitsByScope)

		bumpType, err := n.pkgJsonRepository.GetPackageBumpTypeOutOfCommits(commitsByScope, lib)
		if err != nil {
			fmt.Println("Error getting bump type", err)
			return err
		}

		// bump the version
		version, err := n.BumpNPMPackage(libPath+"/package.json", bumpType)
		if err != nil {
			return err
		}

		changeLogText, err := n.changelogRepository.GetChangelogOutOfCommitMessageAndVersion(commitsByScope, version)
		if err != nil {
			return err
		}

		// If the result shows None ew can get rid of those lines by using this regex expression: ###\s(Deprecated|Added|Bug\sfixes|Breaking Changes)\nNone\.
		fmt.Println("CHANGELOG.md text: ", changeLogText)

		// update the CHANGELOG.md file
		err = n.changelogRepository.UpdateChangeLog(libPath+"/CHANGELOG.md", changeLogText)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NPMPackageService) CreatePR(repoPath string) error {
	err := n.ghRepository.CreatePR(repoPath)
	if err != nil {
		return err
	}

	return nil
}
