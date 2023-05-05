package adapters

import (
	"os"
	"os/exec"
	"path/filepath"
)

type BumpNPMPackage struct {
	BasePath string
}

func NewBumpNPMPackage(basePath string) *BumpNPMPackage {
	return &BumpNPMPackage{
		BasePath: basePath,
	}
}

func (n *BumpNPMPackage) Bump(bumpType string, libPath string) error {
	resolvedPath := filepath.Join(n.BasePath, libPath)
	cmdExec := exec.Command("npm", "version", bumpType, "--no-git-tag-version")
	cmdExec.Dir = filepath.Dir(resolvedPath)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err := cmdExec.Run()

	if err != nil {
		return err
	}
	return nil
}

func (n *BumpNPMPackage) BumpPreRelease(bumpType string, libPath string, preReleaseID string) error {
	resolvedPath := filepath.Join(n.BasePath, libPath)
	cmdExec := exec.Command("npm", "version", bumpType, "--preid", preReleaseID, "--no-git-tag-version")
	cmdExec.Dir = filepath.Dir(resolvedPath)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err := cmdExec.Run()

	if err != nil {
		return err
	}
	return nil
}
