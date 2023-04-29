package git_util

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func GetCommitMessagesByDate(date string, pathToRepo string) (string, error) {
	today := time.Now().Format("2006-01-02T00:00:00")
	cmdDate := today
	if date != "" {
		cmdDate = date
	}

	cmd := exec.Command("git", "log", "--pretty=format:%s%n%b", "--since="+cmdDate, "main")
	cmd.Dir = pathToRepo
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func GetAffectedLibs(pathToRepo string, configuredProjects []string, since string, branch string) ([]string, error) {
	today := time.Now().Format("2006-01-02T00:00:00")
	cmdDate := today

	if since != "" {
		cmdDate = since
	}

	cmd := exec.Command("git", "log", "--pretty=format:", "--name-only", fmt.Sprintf("--since=%s", cmdDate), branch)
	cmd.Dir = pathToRepo
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting affected libs", err)
		return []string{}, err
	}
	fmt.Println("affected libs", string(out))

	affectedLibs := []string{}
	for _, lib := range configuredProjects {
		if strings.Contains(string(out), lib) {
			affectedLibs = append(affectedLibs, lib)
		}
	}

	return affectedLibs, nil
}
