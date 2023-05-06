package domain

import "fmt"

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
