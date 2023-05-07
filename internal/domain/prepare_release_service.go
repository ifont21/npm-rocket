package domain

import (
	"fmt"
)

type PrepareReleaseService struct {
	commitService         CommitsService
	prepareReleasePackage PrepareReleasePackageService
}

func NewPrepareReleaseService(
	commitService CommitsService,
	prepareReleasePackage PrepareReleasePackageService,
) PrepareReleaseService {
	return PrepareReleaseService{
		commitService:         commitService,
		prepareReleasePackage: prepareReleasePackage,
	}
}

func (p PrepareReleaseService) PreReleasePackages() error {
	affectedLibs, err := p.commitService.GetAffectedLibraries("")
	if err != nil {
		fmt.Println("Error getting listAffected", err)
		return err
	}
	fmt.Println("Affected libs: ", affectedLibs)

	commitMessages, err := p.commitService.GetCommitMessagesByDate("")
	if err != nil {
		fmt.Println("Error getting commitMessages", err)
		return err
	}
	fmt.Println("Commit messages: ", commitMessages)

	var processedPackage = make(chan string, len(affectedLibs))

	for _, lib := range affectedLibs {
		go p.prepareReleasePackage.PreparePackage(commitMessages, lib, processedPackage)
	}

	for i := 0; i < len(affectedLibs); i++ {
		fmt.Println(<-processedPackage)
	}

	return nil
}
