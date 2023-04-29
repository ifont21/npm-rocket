package internal

type ChangeLogRepository interface {
	UpdateChangeLog(filePath string, newChangeLog string) error
	GetChangelogOutOfCommitMessageAndVersion(commitMessage string, version string) (string, error)
}
