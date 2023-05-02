package adapters

import (
	"fmt"
	"os/exec"
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

	cmd := exec.Command("git", "log", "--pretty=format:%s%n%b", "--since="+cmdDate, branch)
	cmd.Dir = g.BasePath
	out, err := cmd.Output()
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
