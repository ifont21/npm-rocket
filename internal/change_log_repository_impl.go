package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/ifont21/pre-releaser-cli/internal/files"
	"github.com/ifont21/pre-releaser-cli/internal/gpt"
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

func (c ChangeLogRepositoryImpl) GetChangelogOutOfCommitMessageAndVersion(commitMessage string, version string) (string, error) {
	commit := fmt.Sprintf("commit message:\n%s\n", commitMessage)
	mainPrompt := "can you generate a changelog based on these commits:"
	prompt := fmt.Sprintf("%s\n%s", mainPrompt, commit)

	gptHandler := gpt.NewGPTHandler(os.Getenv("OPENAI_TOKEN"))
	changelogGenerated, err := gptHandler.GetAnswerFromChat(prompt)
	if err != nil {
		return "", err
	}
	changelogGenerated = fmt.Sprintf("## %s\n\n%s", version, changelogGenerated)

	if strings.Contains(changelogGenerated, "Ignoring") {
		changelogGenerated = ""
	}

	return changelogGenerated, nil
}
