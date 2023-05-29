package domain

type BumpPackageJSONService struct {
	bumpPackageJSON   BumpPackageJSON
	fileRepository    FileRepository
	actionSuggestions ActionSuggestions
	config            ConfigRepository
	preRelease        bool
}

func NewBumpPackageJSONService(
	bumpPackageJSON BumpPackageJSON,
	fileRepository FileRepository,
	actionSuggestions ActionSuggestions,
	config ConfigRepository,
	preRelease bool) BumpPackageJSONService {
	return BumpPackageJSONService{
		bumpPackageJSON:   bumpPackageJSON,
		fileRepository:    fileRepository,
		actionSuggestions: actionSuggestions,
		config:            config,
		preRelease:        preRelease,
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

	if b.preRelease {
		preReleaseID, err := b.config.GetPreReleaseID()
		if err != nil {
			return "", err
		}
		preReleaseBumpType := GetBumpTypePreRelease(packageJSON["version"].(string), bumpType)
		err = b.bumpPackageJSON.BumpPreRelease(preReleaseBumpType, libPath, preReleaseID)
		if err != nil {
			return "", err
		}
	} else {
		err = b.bumpPackageJSON.Bump(bumpType, libPath)
		if err != nil {
			return "", err
		}
	}

	newPackageJSON, err := b.fileRepository.GetJSONFileContent(libPath)
	if err != nil {
		return "", err
	}

	return newPackageJSON["version"].(string), nil
}
