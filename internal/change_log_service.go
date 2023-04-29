package internal

type ChangeLogPackageService struct {
	changeLogRepository ChangeLogRepository
}

func NewChangeLogPackageService(changeLogRepository ChangeLogRepository) *ChangeLogPackageService {
	return &ChangeLogPackageService{
		changeLogRepository: changeLogRepository,
	}
}

func (n *ChangeLogPackageService) UpdateChangelog(filePath string, latestChangelog string) error {
	err := n.changeLogRepository.UpdateChangeLog(filePath, latestChangelog)
	if err != nil {
		return err
	}

	return nil
}

func (n *ChangeLogPackageService) GenerateChangelog(commitMessage string, version string) (string, error) {
	changeLogText, err := n.changeLogRepository.GetChangelogOutOfCommitMessageAndVersion(commitMessage, version)
	if err != nil {
		return "", err
	}

	return changeLogText, nil
}
