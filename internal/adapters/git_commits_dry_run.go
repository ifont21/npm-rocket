package adapters

import (
	"os"
	"path/filepath"
)

type GitCommitsDryRun struct {
	config Config
}

func NewGitCommitsDryRun(config Config) *GitCommitsDryRun {
	return &GitCommitsDryRun{
		config: config,
	}
}

func (g *GitCommitsDryRun) GetCommitMessagesByDate(since string, branch string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dryRunCommitsPath, err := g.config.GetDryRunCommitsFilePath()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(filepath.Join(currentDir, dryRunCommitsPath))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (g *GitCommitsDryRun) GetAffectedPaths(since string, branch string) (string, error) {
	return "", nil
}
