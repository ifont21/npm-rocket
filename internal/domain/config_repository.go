package domain

type Package struct {
	Name string
	Path string
}

type Repository struct {
	Owner  string
	Name   string
	Branch string
}

type ConfigRepository interface {
	GetConfiguredLibraries() ([]Package, error)
	GetPreReleaseID() (string, error)
	GetBasedBranch() (string, error)
	GetDryRunCommitsFilePath() (string, error)
	GetRepositoryConfig() (Repository, error)
}
