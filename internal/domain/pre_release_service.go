package domain

import "fmt"

type PreReleaseService struct {
	commitService            CommitsService
	bumpPackageJSONService   BumpPackageJSONService
	generateChangelogService GenerateChangelogService
	fileRepository           FileRepository
}

func NewPreReleaseService(
	commitService CommitsService,
	bumpPackageJSONService BumpPackageJSONService,
	generateChangelogService GenerateChangelogService,
	fileRepository FileRepository) PreReleaseService {
	return PreReleaseService{
		commitService:            commitService,
		bumpPackageJSONService:   bumpPackageJSONService,
		generateChangelogService: generateChangelogService,
		fileRepository:           fileRepository,
	}
}

func (p PreReleaseService) PreReleasePackages() error {
	affectedLibs, err := p.commitService.GetAffectedLibraries("")
	if err != nil {
		fmt.Println("Error getting listAffected", err)
		return err
	}
	fmt.Println("Affected libs: ", affectedLibs)

	commitMessages, err := p.commitService.GetCommitMessagesByDate("", "main")
	if err != nil {
		fmt.Println("Error getting commitMessages", err)
		return err
	}
	fmt.Println("Commit messages: ", commitMessages)

	for _, lib := range affectedLibs {
		/*
		* 1. get the commits by scope.
		 */
		commitsByScope, err := p.commitService.FilterCommitMessageByScope(commitMessages, lib.Name)
		if err != nil {
			fmt.Println("Error getting commits by scope", err)
			return err
		}
		if commitsByScope == "" {
			fmt.Printf("No commits for the scope `%s`\n", lib)
			continue
		}

		fmt.Printf("Text selected by the scope `%s`\n%s", lib, commitsByScope)

		/*
		* 2. bump package.json (libs/<lib>/package.json) by commits. (major, minor, patch) and return the new version.
		 */
		newVersion, err := p.bumpPackageJSONService.BumpPackageByCommits(commitsByScope, fmt.Sprintf("libs/%s/package.json", lib.Path))
		if err != nil {
			fmt.Println("Error bumping package.json", err)
			return err
		}

		/*
		* 3. generate changelog out of the commit message and update it in the file system (libs/<lib>/CHANGELOG.md) with the new version.
		 */
		generatedChangelog, err := p.generateChangelogService.GenerateByCommits(commitsByScope, newVersion)
		if err != nil {
			fmt.Println("Error generating changelog", err)
			return err
		}
		fmt.Println("Generated Changelog: ", generatedChangelog)

		err = p.fileRepository.UpdateTopOfTheFileContent(generatedChangelog, fmt.Sprintf("libs/%s/CHANGELOG.md", lib.Path))
		if err != nil {
			fmt.Println("Error updating changelog", err)
			return err
		}

	}

	return nil
}
