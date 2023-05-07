package adapters

import (
	"github.com/ifont21/pre-releaser-cli/internal/domain"
	"gopkg.in/yaml.v3"
)

type Package struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type PreRelease struct {
	ID string `yaml:"id"`
}

type Commits struct {
	TestFile string `yaml:"test-file"`
}

type PreReleaserYaml struct {
	Repository struct {
		Owner  string `yaml:"owner"`
		Name   string `yaml:"name"`
		Branch string `yaml:"branch"`
	} `yaml:"repository"`
	Commits    Commits    `yaml:"commits"`
	PreRelease PreRelease `yaml:"pre-release"`
	Libs       []Package  `yaml:"libs"`
}

func castDomainPackage(yamlPackages []Package) []domain.Package {
	var domainPackages []domain.Package
	for _, yamlPackage := range yamlPackages {
		domainPackages = append(domainPackages, domain.Package{
			Name: yamlPackage.Name,
			Path: yamlPackage.Path,
		})
	}
	return domainPackages
}

type Config struct {
	fileRepository FileRepository
}

func NewConfig(fileRepository FileRepository) Config {
	return Config{
		fileRepository: fileRepository,
	}
}

func (c Config) GetConfiguredLibraries() ([]domain.Package, error) {
	var preReleaseConfig PreReleaserYaml
	config, err := c.fileRepository.GetPlainFileContent("releaser.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(config, &preReleaseConfig)
	if err != nil {
		return nil, err
	}

	return castDomainPackage(preReleaseConfig.Libs), nil
}

func (c Config) GetPreReleaseID() (string, error) {
	var preReleaseConfig PreReleaserYaml
	config, err := c.fileRepository.GetPlainFileContent("releaser.yaml")
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(config, &preReleaseConfig)
	if err != nil {
		return "", err
	}

	return preReleaseConfig.PreRelease.ID, nil
}

func (c Config) GetBasedBranch() (string, error) {
	var preReleaseConfig PreReleaserYaml
	config, err := c.fileRepository.GetPlainFileContent("releaser.yaml")
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(config, &preReleaseConfig)
	if err != nil {
		return "", err
	}

	return preReleaseConfig.Repository.Branch, nil
}

func (c Config) GetDryRunCommitsFilePath() (string, error) {
	var preReleaseConfig PreReleaserYaml
	config, err := c.fileRepository.GetPlainFileContent("releaser.yaml")
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(config, &preReleaseConfig)
	if err != nil {
		return "", err
	}

	return preReleaseConfig.Commits.TestFile, nil
}
