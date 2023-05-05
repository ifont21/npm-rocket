package domain

type BumpPackageJSON interface {
	Bump(bumpType string, libPath string) error
	BumpPreRelease(bumpType string, libPath string, preReleaseID string) error
}
