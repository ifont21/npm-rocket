package domain

type Package struct {
	Name string
	Path string
}

type PreReleaserYaml struct {
	Repository struct {
		Owner string
		Name  string
	}
	Libs []Package
}

type ConfigRepository interface {
	GetConfiguredLibraries() ([]Package, error)
}
