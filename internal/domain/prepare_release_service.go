package domain

import (
	"fmt"
	"time"
)

func getReleaseBranchName(isPreRelease bool) string {
	today := time.Now().Format("2006-01-02")

	if isPreRelease {
		return fmt.Sprintf("pkg-pre-release-%s", today)
	}
	return fmt.Sprintf("pkg-release-%s", today)
}

func getTitlePR(isPreRelease bool) string {
	if isPreRelease {
		return fmt.Sprintf("Packages Pre-Release %s", time.Now().Format("2006-01-02"))
	}

	return fmt.Sprintf("Packages Stable Release %s", time.Now().Format("2006-01-02"))
}

type PrepareReleaseService struct {
	commitService         CommitsService
	prepareReleasePackage PrepareReleasePackageService
	gitChangesRepository  GitChangesRepository
	prRepository          PRRepository
}

func NewPrepareReleaseService(
	commitService CommitsService,
	prepareReleasePackage PrepareReleasePackageService,
	gitChanges GitChangesRepository,
	prRepository PRRepository,
) PrepareReleaseService {
	return PrepareReleaseService{
		commitService:         commitService,
		prepareReleasePackage: prepareReleasePackage,
		gitChangesRepository:  gitChanges,
		prRepository:          prRepository,
	}
}

func (p PrepareReleaseService) PreReleasePackages(isPreRelease bool, noCommit bool) error {
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

	releaseBranchName := getReleaseBranchName(isPreRelease)

	// If noCommit is true, we don't want to commit and push changes to the remote
	// If true we don't want to create a PR either
	if noCommit {
		return nil
	}

	err = p.gitChangesRepository.AddAndCommitChanges(releaseBranchName)
	if err != nil {
		fmt.Println("Error on add and commit changes", err)
		return err
	}

	affectedPackages := []string{}
	for _, lib := range affectedLibs {
		affectedPackages = append(affectedPackages, lib.Name)
	}

	err = p.prRepository.CreatePR(PR{
		TitlePR:    getTitlePR(isPreRelease),
		BranchName: releaseBranchName,
	}, affectedPackages)
	if err != nil {
		fmt.Println("Error on create PR", err)
		return err
	}

	return nil
}
