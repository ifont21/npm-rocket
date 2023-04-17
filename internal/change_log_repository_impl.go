package internal

import (
	"fmt"

	"github.com/ifont21/pre-releaser-cli/internal/files"
)

type ChangeLogRepositoryImpl struct{}

func NewChangeLogRepositoryImpl() *ChangeLogRepositoryImpl {
	return &ChangeLogRepositoryImpl{}
}

func (c ChangeLogRepositoryImpl) UpdateChangeLog(filePath string, newChangeLog string) error {
	err := files.AddTextOnTopOfFile(filePath, newChangeLog)
	if err != nil {
		return err
	}

	return nil
}

func (c ChangeLogRepositoryImpl) GetChangeLogOutOfCommitMessageAndVersion(commitMessage string, version string) (string, error) {
	changelogTemplate := fmt.Sprintf("## %s\n\n%s\n\n", version, commitMessage)

	return changelogTemplate, nil
}
