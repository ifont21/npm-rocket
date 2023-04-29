package files

type PreReleaserYaml struct {
	Repository struct {
		Owner string `yaml:"owner"`
		Name  string `yaml:"name"`
	} `yaml:"repository"`
	Libs []string `yaml:"libs"`
}
