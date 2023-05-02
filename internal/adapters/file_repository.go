package adapters

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileRepository struct {
	basePath string
}

// NewFileRepository creates a new instance of FileRepository
func NewFileRepository(basePath string) FileRepository {
	return FileRepository{
		basePath: basePath,
	}
}

// get the content of a plain text file
func (f FileRepository) GetPlainFileContent(libPath string) ([]byte, error) {
	pathResolved := fmt.Sprintf("%s/%s", f.basePath, libPath)
	content, err := os.ReadFile(pathResolved)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

// read a json file and return the content as a map
func (f FileRepository) GetJSONFileContent(libPath string) (map[string]interface{}, error) {
	pathResolved := fmt.Sprintf("%s/%s", f.basePath, libPath)
	var jsonFile map[string]interface{}
	content, err := os.ReadFile(pathResolved)
	if err != nil {
		return map[string]interface{}{}, err
	}

	err = json.Unmarshal(content, &jsonFile)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return jsonFile, nil
}

// add text on top of a file
func (f FileRepository) UpdateTopOfTheFileContent(newContent string, libPath string) error {
	pathResolved := fmt.Sprintf("%s/%s", f.basePath, libPath)
	content, err := os.ReadFile(pathResolved)
	if err != nil {
		return err
	}

	contentStr := string(content)
	contentStr = fmt.Sprintf("%s\n\n%s", newContent, contentStr)

	err = os.WriteFile(pathResolved, []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
