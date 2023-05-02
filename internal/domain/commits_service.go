package domain

import "regexp"

type CommitsService struct {
	actionSuggestions   ActionSuggestions
	gitCommitRepository GitCommitsRepository
	config              ConfigRepository
}

func NewCommitsService(actionSuggestions ActionSuggestions, gitCommitRepository GitCommitsRepository, config ConfigRepository) CommitsService {
	return CommitsService{
		actionSuggestions:   actionSuggestions,
		gitCommitRepository: gitCommitRepository,
		config:              config,
	}
}

func (c CommitsService) FilterCommitMessageByScope(commits string, scope string) (string, error) {
	commits, err := c.actionSuggestions.GetFilteredCommitsByScope(commits, scope)
	if err != nil {
		return "", err
	}
	// evaluate regex to commits
	// return commits
	regex, err := regexp.Compile("(no commits|no commits related.*)")
	if err != nil {
		return "", err
	}
	if regex.MatchString(commits) {
		return "", nil
	}

	return commits, nil
}

func (c CommitsService) GetCommitMessagesByDate(since string, branch string) (string, error) {
	commits, err := c.gitCommitRepository.GetCommitMessagesByDate(since, branch)
	if err != nil {
		return "", err
	}

	return commits, nil
}

func (c CommitsService) GetAffectedLibraries(base string) ([]Package, error) {
	configuredProjects, err := c.config.GetConfiguredLibraries()
	if err != nil {
		return nil, err
	}

	return configuredProjects, nil
}
