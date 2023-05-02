package domain

import (
	"fmt"
	"strings"
)

type GenerateChangelogService struct {
	ActionSuggestions ActionSuggestions
}

func NewGenerateChangelogService(actionSuggestions ActionSuggestions) GenerateChangelogService {
	return GenerateChangelogService{
		ActionSuggestions: actionSuggestions,
	}
}

func (g GenerateChangelogService) GenerateByCommits(commits string, newVersion string) (string, error) {
	suggestedText, err := g.ActionSuggestions.GetSuggestedChangelogOutOfCommits(commits)
	if err != nil {
		return "", err
	}
	changelogGenerated := fmt.Sprintf("## %s\n\n%s", newVersion, suggestedText)

	if strings.Contains(changelogGenerated, "Ignoring") {
		changelogGenerated = ""
	}

	return changelogGenerated, nil
}
