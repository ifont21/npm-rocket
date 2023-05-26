package adapters

import (
	"github.com/ifont21/pre-releaser-cli/internal/helpers"
)

type GitChanges struct {
	BasePath string
}

func NewGitChanges(basePath string) *GitChanges {
	return &GitChanges{
		BasePath: basePath,
	}
}

func (g *GitChanges) AddAndCommitChanges(branchName string) error {
	// verify that exist changes to commit
	err := helpers.CheckStatus(g.BasePath)
	if err != nil {
		return err
	}

	err = helpers.CreateBranch(g.BasePath, branchName)
	if err != nil {
		return err
	}

	err = helpers.GitAddChanges(g.BasePath)
	if err != nil {
		return err
	}

	err = helpers.GitCommitChanges(g.BasePath, "chore: release "+branchName)
	if err != nil {
		return err
	}

	err = helpers.GitPushChanges(g.BasePath, branchName)
	if err != nil {
		return err
	}

	return nil
}
