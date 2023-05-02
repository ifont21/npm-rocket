package domain

import "strings"

func GetAffectedProjects(configuredProjects []string, pathsInCommitMessages string) []string {
	affectedLibs := []string{}
	for _, lib := range configuredProjects {
		if strings.Contains(string(pathsInCommitMessages), lib) {
			affectedLibs = append(affectedLibs, lib)
		}
	}

	return affectedLibs
}
