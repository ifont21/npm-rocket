package adapters

import (
	"github.com/ifont21/pre-releaser-cli/internal/domain"
	"github.com/ifont21/pre-releaser-cli/internal/stubs"
)

func NewPreReleaserContainer(basePath string, openAIToken string) domain.PreReleaseService {
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
	bumpPackageJSONService := domain.NewBumpPackageJSONService(bumpPackageJSON, fileRepository, suggestions)
	generateChangelogService := domain.NewGenerateChangelogService(suggestions)

	return domain.NewPreReleaseService(commitService, bumpPackageJSONService, generateChangelogService, fileRepository)
}
