package internal

import (
	"fmt"
	"os"

	"github.com/ifont21/pre-releaser-cli/internal/files"
	"github.com/ifont21/pre-releaser-cli/internal/git_util"
	"github.com/ifont21/pre-releaser-cli/internal/gpt"
)

type GitRepositoryImpl struct{}

func NewGitRepositoryImpl() *GitRepositoryImpl {
	return &GitRepositoryImpl{}
}

func (g *GitRepositoryImpl) SummarizeCommitsByScope(scope string, commits string) (string, error) {
	prompt := fmt.Sprintf("Select only the text related to `%s`\n%s", scope, commits)

	token := os.Getenv("OPENAI_TOKEN")
	gptHandler := gpt.NewGPTHandler(token)
	filteredCommit, err := gptHandler.GetAnswerFromChat(prompt)
	if err != nil {
		return "", err
	}

	return filteredCommit, nil
}

func (g *GitRepositoryImpl) GetAffectedLibsFromGitCommits(repoPath string, base string) ([]string, error) {
	releaserConfigFilePath := repoPath + "/pre-releaser.yaml"
	configureLibs, err := files.GetReleaseCLIConfig(releaserConfigFilePath)
	if err != nil {
		fmt.Println("Error getting configureLibs", err)
		return []string{}, err
	}
	baseConfig := base
	if base == "" {
		baseConfig = "main"
	}

	listAffected, err := git_util.GetAffectedLibs(repoPath, configureLibs.Libs, "", baseConfig)
	if err != nil {
		fmt.Println("Error getting listAffected", err)
		return []string{}, err
	}

	return listAffected, nil
}
