package domain

import "fmt"

func GetPackageProperty(packageJSON map[string]interface{}, property string) (interface{}, error) {
	if packageJSON[property] == nil {
		return nil, fmt.Errorf("property %s not found", property)
	}

	return packageJSON[property], nil
}
