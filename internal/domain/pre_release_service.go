package domain

import (
	"fmt"
)

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

	var processedPackage = make(chan string, len(affectedLibs))

	for _, lib := range affectedLibs {
		go func(lib Package) {
			commitsByScope, err := p.commitService.FilterCommitMessageByScope(commitMessages, lib.Name)
			if err != nil {
				processedPackage <- fmt.Sprintf("Channel: Error getting commits by scope `%s`\n", lib.Name)
				return
			}
			if commitsByScope == "" {
				processedPackage <- fmt.Sprintf("Channel: No commits for the scope `%s`\n", lib.Name)
				return
			}

			newVersion, err := p.bumpPackageJSONService.BumpPackageByCommits(commitsByScope, fmt.Sprintf("libs/%s/package.json", lib.Path))
			if err != nil {
				processedPackage <- fmt.Sprintf("Channel: Error bumping package.json `%s`\n", lib.Name)
				return
			}

			generatedChangelog, err := p.generateChangelogService.GenerateByCommits(commitsByScope, newVersion)
			if err != nil {
				processedPackage <- fmt.Sprintf("Channel: Error generating the CHANGELOG for `%s`\n", lib.Name)
				return
			}

			err = p.fileRepository.UpdateTopOfTheFileContent(generatedChangelog, fmt.Sprintf("libs/%s/CHANGELOG.md", lib.Path))
			if err != nil {
				processedPackage <- fmt.Sprintf("Channel: Error Updating the CHANGELOG for `%s`\n", lib.Name)
				return
			}
			processedPackage <- fmt.Sprintf("New version generated for `%s`: ****************************************\n%s\n", lib.Name, generatedChangelog)
		}(lib)
	}

	for i := 0; i < len(affectedLibs); i++ {
		fmt.Println(<-processedPackage)
	}

	/* for _, lib := range affectedLibs {
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

		newVersion, err := p.bumpPackageJSONService.BumpPackageByCommits(commitsByScope, fmt.Sprintf("libs/%s/package.json", lib.Path))
		if err != nil {
			fmt.Println("Error bumping package.json", err)
			return err
		}

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

	} */

	return nil
}
