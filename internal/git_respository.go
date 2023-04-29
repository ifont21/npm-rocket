package internal

type GitRepository interface {
	SummarizeCommitsByScope(scope string, commits string) (string, error)
	GetAffectedLibsFromGitCommits(repoPath string, base string) ([]string, error)
}
