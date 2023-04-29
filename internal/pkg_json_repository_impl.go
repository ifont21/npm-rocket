package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ifont21/pre-releaser-cli/internal/files"
	"github.com/ifont21/pre-releaser-cli/internal/gpt"
)

type PkgJSONRepositoryImpl struct{}

func NewPkgJSONRepositoryImpl() *PkgJSONRepositoryImpl {
	return &PkgJSONRepositoryImpl{}
}

func (p *PkgJSONRepositoryImpl) GetPackageJSON(filePath string) (PackageJSON, error) {
	var packageJSON PackageJSON
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JSON filePath content %s\n", string(file))
	json.Unmarshal(file, &packageJSON)

	return packageJSON, nil
}

func (p *PkgJSONRepositoryImpl) BumpNPMPackage(filePath string, bumpType string) (string, error) {
	packageJSON, err := p.GetPackageJSON(filePath)
	if err != nil {
		return "", err
	}

	fmt.Println("Package name", packageJSON["name"])
	fmt.Println("Initial package version", packageJSON["version"])
	fmt.Printf("Bump type: %s\n", bumpType)
	fmt.Println("Bumping version...")

	cmdExec := exec.Command("npm", "version", bumpType, "--no-git-tag-version")
	cmdExec.Dir = filepath.Dir(filePath)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err = cmdExec.Run()

	if err != nil {
		return "", err
	}

	version, err := files.GetJSONPropertyFromFile(filePath, "version")
	if err != nil {
		return "", err
	}

	return version, nil
}

func (p *PkgJSONRepositoryImpl) GetPackageBumpTypeOutOfCommits(commits string, lib string) (string, error) {

	request := `based on the commit message how you suggest to bump the library
		- major\n
		- minor\n
		- patch
	`
	prompt := fmt.Sprintf("%s\n%s", request, commits)

	gptHandler := gpt.NewGPTHandler(os.Getenv("OPENAI_TOKEN"))
	text, err := gptHandler.GetAnswerFromChat(prompt)
	if err != nil {
		return "", err
	}
	bumpType := "patch"

	if strings.Contains(text, "major") {
		bumpType = "major"
	}

	if strings.Contains(text, "minor") {
		bumpType = "minor"
	}
	fmt.Println("Bump Type suggested -->: ", bumpType)

	return bumpType, nil
}
