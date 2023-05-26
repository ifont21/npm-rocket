package helpers

import (
	"errors"
	"fmt"
	"os/exec"
)

func CheckStatus(dirPath string) error {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error on getting git status")
		return err
	}
	if len(out) == 0 {
		return errors.New("no changes to commit for this release")
	}

	return nil
}

func CreateBranch(dirPath string, branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dirPath
	_, err := cmd.Output()
	if err != nil {
		return errors.New("error on creating new branch")
	}

	return nil
}

func GitAddChanges(dirPath string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dirPath
	_, err := cmd.Output()
	if err != nil {
		return errors.New("error on adding changes")
	}

	return nil
}

func GitCommitChanges(dirPath string, commitMessage string) error {
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Dir = dirPath
	_, err := cmd.Output()
	if err != nil {
		return errors.New("error on committing changes")
	}

	return nil
}

func GitPushChanges(dirPath string, branchName string) error {
	cmd := exec.Command("git", "push", "-u", "origin", branchName)
	cmd.Dir = dirPath
	_, err := cmd.Output()
	if err != nil {
		return errors.New("error on pushing changes")
	}

	return nil
}
