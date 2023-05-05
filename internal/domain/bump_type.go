package domain

import (
	"strings"
)

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

func tearDownVersion(version string) (string, string, string) {
	stableVersion := version
	if strings.Contains(version, "-") {
		stableVersion = strings.Split(version, "-")[0]
	}

	versionTokens := strings.Split(stableVersion, ".")
	return versionTokens[0], versionTokens[1], versionTokens[2]
}

func preReleaseMajor(currentVersion string, isPreReleased bool) string {
	_, minor, patch := tearDownVersion(currentVersion)
	if isPreReleased && patch == "0" && minor == "0" {
		return "prerelease"
	}
	return "premajor"
}

func preReleaseMinor(currentVersion string, isPreReleased bool) string {
	_, _, patch := tearDownVersion(currentVersion)
	if isPreReleased && patch == "0" {
		return "prerelease"
	}
	return "preminor"
}

func preReleasePatch(isPreReleased bool) string {
	if isPreReleased {
		return "prerelease"
	}
	return "prepatch"
}

func GetBumpTypePreRelease(currentVersion string, suggestedType string) string {
	isPreReleased := strings.Contains(currentVersion, "-")

	if suggestedType == "major" {
		return preReleaseMajor(currentVersion, isPreReleased)
	}

	if suggestedType == "minor" {
		return preReleaseMinor(currentVersion, isPreReleased)
	}

	if suggestedType == "patch" {
		return preReleasePatch(isPreReleased)
	}

	return "prerelease"
}
