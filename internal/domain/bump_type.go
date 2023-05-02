package domain

import "strings"

func GetBumpTypeOutOfText(text string) string {
	bumpType := "patch"

	if strings.Contains(text, "major") {
		bumpType = "major"
	}

	if strings.Contains(text, "minor") {
		bumpType = "minor"
	}

	return bumpType
}
