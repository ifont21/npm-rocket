package domain

type ActionSuggestions interface {
	GetSuggestedChangelogOutOfCommits(commits string) (string, error)
	GetBumpTypeSuggestionOutOfCommits(commits string) (string, error)
	GetFilteredCommitsByScope(commits string, scope string) (string, error)
}
