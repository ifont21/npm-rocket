package stubs

import (
	"os"
	"path/filepath"
)

type GitCommitsStub struct{}

func NewGitCommitsStub() *GitCommitsStub {
	return &GitCommitsStub{}
}

func (g *GitCommitsStub) GetCommitMessagesByDate(since string, branch string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(filepath.Join(currentDir, "resources", "commits.txt"))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (g *GitCommitsStub) GetAffectedPaths(since string, branch string) (string, error) {
	return "", nil
}
