package internal

type PackageJSON = map[string]interface{}

type PkgJSONRepository interface {
	GetPackageJSON(filePath string) (PackageJSON, error)
	BumpNPMPackage(filePath string, bumpType string) (string, error)
}
