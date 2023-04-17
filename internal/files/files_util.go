package files

import (
	"encoding/json"
	"os"
)

func AddTextOnTopOfFile(filePath string, text string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	contentStr = text + contentStr

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
