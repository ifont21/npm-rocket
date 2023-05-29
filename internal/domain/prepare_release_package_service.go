package domain

import (
	"fmt"
	"sync"
)

type PreChangelog struct {
	commitMessages string
	packageVersion string
	packageInfo    Package
}

func NewPreChangelog(messages, version string, packageInfo Package) PreChangelog {
	return PreChangelog{
		commitMessages: messages,
		packageVersion: version,
		packageInfo:    packageInfo,
	}
}

type PrepareReleasePackageService struct {
	commitService            CommitsService
	bumpPackageJSONService   BumpPackageJSONService
	generateChangelogService GenerateChangelogService
	fileRepository           FileRepository
}

func NewPrepareReleasePackageService(commitService CommitsService,
	bumpPackageService BumpPackageJSONService,
	generateChangelogService GenerateChangelogService,
	fileRepository FileRepository,
) PrepareReleasePackageService {
	return PrepareReleasePackageService{
		commitService:            commitService,
		bumpPackageJSONService:   bumpPackageService,
		generateChangelogService: generateChangelogService,
		fileRepository:           fileRepository,
	}
}

func (r PrepareReleasePackageService) BumpPackage(commitMessages string, packageLibrary Package, wg *sync.WaitGroup, preChangelog chan<- PreChangelog) {
	defer wg.Done()
	fmt.Printf("Started processing package: %s\n", packageLibrary.Name)

	commitsByScope, err := r.commitService.FilterCommitMessageByScope(commitMessages, packageLibrary.Name)
	if err != nil {
		fmt.Println("Error getting commits by scope", err)
		preChangelog <- PreChangelog{}
		return
	}

	if commitsByScope == "" {
		fmt.Println("No commits for package", packageLibrary.Name)
		preChangelog <- PreChangelog{}
		return
	}

	fmt.Printf("Bumping package %s\n", packageLibrary.Name)
	newVersion, err := r.bumpPackageJSONService.BumpPackageByCommits(commitsByScope, fmt.Sprintf("libs/%s/package.json", packageLibrary.Path))
	if err != nil {
		fmt.Println("Error bumping package", err)
		preChangelog <- PreChangelog{}
		return
	}

	preChangelog <- NewPreChangelog(commitsByScope, newVersion, packageLibrary)
	fmt.Printf("Bumping package %s is completed \n", packageLibrary.Name)
}

func (r PrepareReleasePackageService) GenerateChangelogForPackage(preChangelog PreChangelog, wg *sync.WaitGroup, generated chan<- string) {
	defer wg.Done()
	fmt.Printf("Started generating changelog for package: %s\n", preChangelog.packageInfo.Name)

	if preChangelog.packageVersion == "" {
		fmt.Printf("Skipping changelog generation for `%s`\n", preChangelog.packageInfo.Name)
		generated <- ""
		return
	}

	generatedChangelog, err := r.generateChangelogService.GenerateByCommits(preChangelog.commitMessages, preChangelog.packageVersion)
	if err != nil {
		fmt.Printf("Channel: Error generating the CHANGELOG for `%s`\n", preChangelog.packageInfo.Name)
		generated <- ""
		return
	}

	err = r.fileRepository.UpdateTopOfTheFileContent(generatedChangelog, fmt.Sprintf("libs/%s/CHANGELOG.md", preChangelog.packageInfo.Path))
	if err != nil {
		fmt.Printf("Channel: Error Updating the CHANGELOG for `%s`\n", preChangelog.packageInfo.Name)
		generated <- ""
		return
	}

	generated <- fmt.Sprintf("New version generated for `%s`: ****************************************\n%s\n", preChangelog.packageInfo.Name, generatedChangelog)
	fmt.Printf("Changelog generated for %s package\n", preChangelog.packageInfo.Name)
}

func (r PrepareReleasePackageService) PreparePackage(commitMessages string, packageLibrary Package, processedPackage chan<- string) {
	commitsByScope, err := r.commitService.FilterCommitMessageByScope(commitMessages, packageLibrary.Name)
	if err != nil {
		processedPackage <- fmt.Sprintf("Channel: Error getting commits by scope `%s`\n", packageLibrary.Name)
		return
	}
	if commitsByScope == "" {
		processedPackage <- fmt.Sprintf("Channel: No commits for the scope `%s`\n", packageLibrary.Name)
		return
	}

	newVersion, err := r.bumpPackageJSONService.BumpPackageByCommits(commitsByScope, fmt.Sprintf("libs/%s/package.json", packageLibrary.Path))
	if err != nil {
		processedPackage <- fmt.Sprintf("Channel: Error bumping package.json `%s`\n", packageLibrary.Name)
		return
	}

	generatedChangelog, err := r.generateChangelogService.GenerateByCommits(commitsByScope, newVersion)
	if err != nil {
		processedPackage <- fmt.Sprintf("Channel: Error generating the CHANGELOG for `%s`\n", packageLibrary.Name)
		return
	}

	err = r.fileRepository.UpdateTopOfTheFileContent(generatedChangelog, fmt.Sprintf("libs/%s/CHANGELOG.md", packageLibrary.Path))
	if err != nil {
		processedPackage <- fmt.Sprintf("Channel: Error Updating the CHANGELOG for `%s`\n", packageLibrary.Name)
		return
	}
	processedPackage <- fmt.Sprintf("New version generated for `%s`: ****************************************\n%s\n", packageLibrary.Name, generatedChangelog)
}
