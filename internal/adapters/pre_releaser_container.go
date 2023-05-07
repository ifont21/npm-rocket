package adapters

import (
	"github.com/ifont21/pre-releaser-cli/internal/domain"
)

func NewPreReleaserContainer(basePath string, openAIToken string, preRelease bool, dryRun bool) domain.PrepareReleaseService {
	// Repositories
	fileRepository := NewFileRepository(basePath)
	// configuration from file
	config := NewConfig(fileRepository)
	bumpPackageJSON := NewBumpNPMPackage(basePath)
	suggestions := NewGPTSuggestion(openAIToken)

	var gitCommitsRepository domain.GitCommitsRepository
	if dryRun {
		gitCommitsRepository = NewGitCommitsDryRun(config)
	} else {
		gitCommitsRepository = NewGitCommits(basePath)
	}

	// Services
	commitService := domain.NewCommitsService(suggestions, gitCommitsRepository, config)
	bumpPackageJSONService := domain.NewBumpPackageJSONService(bumpPackageJSON,
		fileRepository,
		suggestions,
		config,
		preRelease,
	)
	generateChangelogService := domain.NewGenerateChangelogService(suggestions)
	prepareReleasePackageService := domain.NewPrepareReleasePackageService(commitService, bumpPackageJSONService, generateChangelogService, fileRepository)

	return domain.NewPrepareReleaseService(commitService, prepareReleasePackageService)
}
