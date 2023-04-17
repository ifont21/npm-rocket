package internal

import "fmt"

type NPMPackageService struct {
	pkgJsonRepository   PkgJSONRepository
	changelogRepository ChangeLogRepository
}

func NewNPMPackageService(pkgJsonRepository PkgJSONRepository, changelogRepository ChangeLogRepository) *NPMPackageService {
	return &NPMPackageService{
		pkgJsonRepository:   pkgJsonRepository,
		changelogRepository: changelogRepository,
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
func (n *NPMPackageService) BumpNPMPackagesAndChangelog(libPath string, libs []string, commits string) error {
	// loop through the libs and bump the versions
	for _, lib := range libs {
		libPath := libPath + "/" + lib
		// bump the version
		version, err := n.BumpNPMPackage(libPath+"/package.json", "patch")
		if err != nil {
			return err
		}
		fmt.Printf("Bumped version to %s\n", version)

		changeLogText, err := n.changelogRepository.GetChangeLogOutOfCommitMessageAndVersion(commits, version)
		if err != nil {
			return err
		}
		fmt.Printf("generated CHANGELOG is  %s\n", changeLogText)

		// update the CHANGELOG.md file
		err = n.changelogRepository.UpdateChangeLog(libPath+"/CHANGELOG.md", changeLogText)
		if err != nil {
			return err
		}
	}

	return nil
}
