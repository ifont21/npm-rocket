package adapters

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type GitCommits struct {
	BasePath string
}

func NewGitCommits(basePath string) *GitCommits {
	return &GitCommits{
		BasePath: basePath,
	}
}

func (g *GitCommits) GetCommitMessagesByDate(since string, branch string) (string, error) {
	today := time.Now().Format("2006-01-02T00:00:00")
	cmdDate := today
	if since != "" {
		cmdDate = since
	}

	// 1. validate if the current branch is the same as the one passed as parameter
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.BasePath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(string(out)) != branch {
		return "", errors.New("wrong branch, please checkout to " + branch + " branch")
	}

	// 2. pull the latest changes
	fmt.Printf("Pulling latest changes from %s branch...\n", branch)
	cmd = exec.Command("git", "pull", "origin", branch, "--rebase")
	cmd.Dir = g.BasePath
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	cmd = exec.Command("git", "log", "--pretty=format:%s%n%b", "--since="+cmdDate, branch)
	cmd.Dir = g.BasePath
	out, err = cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (g *GitCommits) GetAffectedPaths(since string, branch string) (string, error) {
	today := time.Now().Format("2006-01-02T00:00:00")
	cmdDate := today
	if since != "" {
		cmdDate = since
	}

	cmd := exec.Command("git", "log", "--pretty=format:", "--name-only", fmt.Sprintf("--since=%s", cmdDate), branch)
	cmd.Dir = g.BasePath
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting affected libs", err)
		return "", err
	}
	fmt.Println("out:", string(out))

	return string(out), nil
}
