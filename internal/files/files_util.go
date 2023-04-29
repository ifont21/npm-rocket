package files

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func AddTextOnTopOfFile(filePath string, text string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	contentStr = fmt.Sprintf("%s\n\n%s", text, contentStr)

	err = os.WriteFile(filePath, []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetJSONPropertyFromFile(filePath string, property string) (string, error) {
	var jsonFile map[string]interface{}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(content, &jsonFile)
	if err != nil {
		return "", err
	}

	return jsonFile[property].(string), nil
}

// read a yml file and return the content as a string
func GetReleaseCLIConfig(filePath string) (PreReleaserYaml, error) {
	// read the file
	var preReleaseConfig PreReleaserYaml
	content, err := os.ReadFile(filePath)
	if err != nil {
		return PreReleaserYaml{}, err
	}

	err = yaml.Unmarshal(content, &preReleaseConfig)
	if err != nil {
		return PreReleaserYaml{}, err
	}

	return preReleaseConfig, nil
}
