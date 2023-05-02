package domain

type GitCommitsRepository interface {
	GetCommitMessagesByDate(since string, branch string) (string, error)
	GetAffectedPaths(since string, branch string) (string, error)
}
