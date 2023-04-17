package internal

type ChangeLogRepository interface {
	UpdateChangeLog(filePath string, newChangeLog string) error
	GetChangeLogOutOfCommitMessageAndVersion(commitMessage string, version string) (string, error)
}
