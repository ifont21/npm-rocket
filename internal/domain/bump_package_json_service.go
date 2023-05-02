package domain

import "fmt"

type BumpPackageJSONService struct {
	bumpPackageJSON   BumpPackageJSON
	fileRepository    FileRepository
	actionSuggestions ActionSuggestions
}

func NewBumpPackageJSONService(
	bumpPackageJSON BumpPackageJSON,
	fileRepository FileRepository,
	actionSuggestions ActionSuggestions) BumpPackageJSONService {
	return BumpPackageJSONService{
		bumpPackageJSON:   bumpPackageJSON,
		fileRepository:    fileRepository,
		actionSuggestions: actionSuggestions,
	}
}

func (b BumpPackageJSONService) BumpPackageByCommits(commits string, libPath string) (string, error) {
	bumpTypeSuggestion, err := b.actionSuggestions.GetBumpTypeSuggestionOutOfCommits(commits)
	if err != nil {
		return "", err
	}

	bumpType := GetBumpTypeOutOfText(bumpTypeSuggestion)
	packageJSON, err := b.fileRepository.GetJSONFileContent(libPath)
	if err != nil {
		return "", err
	}
	fmt.Printf("The suggested semantic bump type version for %s is %s\n", packageJSON["name"], bumpType)
	fmt.Printf("Starting to bump NPM package %s\n", packageJSON["name"])
	fmt.Printf("Current version: %s\n", packageJSON["version"])

	err = b.bumpPackageJSON.Bump(bumpType, libPath)
	if err != nil {
		return "", err
	}

	newPackageJSON, err := b.fileRepository.GetJSONFileContent(libPath)
	if err != nil {
		return "", err
	}
	fmt.Printf("New version: %s\n", newPackageJSON["version"])

	return newPackageJSON["version"].(string), nil
}
