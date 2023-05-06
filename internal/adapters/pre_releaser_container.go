package adapters

import (
	"github.com/ifont21/pre-releaser-cli/internal/domain"
	"github.com/ifont21/pre-releaser-cli/internal/stubs"
)

func NewPreReleaserContainer(basePath string, openAIToken string, preRelease bool) domain.PrepareReleaseService {
	// Repositories
	fileRepository := NewFileRepository(basePath)
	// gitCommitsRepository := NewGitCommits(basePath)
	gitCommitsRepository := stubs.NewGitCommitsStub()
	bumpPackageJSON := NewBumpNPMPackage(basePath)
	suggestions := NewGPTSuggestion(openAIToken)
	// configuration from file
	config := NewConfig(fileRepository)
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
