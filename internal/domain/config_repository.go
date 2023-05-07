package domain

type Package struct {
	Name string
	Path string
}

type ConfigRepository interface {
	GetConfiguredLibraries() ([]Package, error)
	GetPreReleaseID() (string, error)
	GetBasedBranch() (string, error)
	GetDryRunCommitsFilePath() (string, error)
}
