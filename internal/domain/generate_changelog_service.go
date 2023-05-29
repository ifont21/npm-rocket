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
	changelogGenerated, err := g.ActionSuggestions.GetSuggestedChangelogOutOfCommits(commits)
	if err != nil {
		return "", err
	}

	beginIndex := strings.Index(changelogGenerated, "- begin")
	endIndex := strings.Index(changelogGenerated, "- end")
	if beginIndex == -1 || endIndex == -1 {
		return "", fmt.Errorf("could not find begin or end in the generated changelog")
	}
	changelogGenerated = changelogGenerated[beginIndex:endIndex]
	changelogGenerated = strings.ReplaceAll(changelogGenerated, "- begin", "")
	changelogGenerated = strings.ReplaceAll(changelogGenerated, "- end", "")
	changelogGenerated = strings.TrimSpace(changelogGenerated)

	changelogGenerated = fmt.Sprintf("## %s\n\n%s", newVersion, changelogGenerated)

	return changelogGenerated, nil
}
