package domain

type GitChangesRepository interface {
	AddAndCommitChanges(branchName string) error
}
