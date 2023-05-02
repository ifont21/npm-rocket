package domain

type BumpPackageJSON interface {
	Bump(bumpType string, libPath string) error
}
