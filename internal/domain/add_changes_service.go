package domain

type AddChangesService struct {
	gitChangesRepository GitChangesRepository
	prRepository         PRRepository
}

func NewAddChangesService(gitChanges GitChangesRepository, prRepo PRRepository) AddChangesService {
	return AddChangesService{
		gitChangesRepository: gitChanges,
		prRepository:         prRepo,
	}
}

func (a AddChangesService) AddChanges() error {
	err := a.gitChangesRepository.AddAndCommitChanges("release")
	if err != nil {
		return err
	}

	err = a.prRepository.CreatePR(PR{
		TitlePR:    "Release-PR-2023-05-22",
		BranchName: "release",
	}, []string{"feature-a-lib", "feature-b-lib"})
	if err != nil {
		return err
	}

	return nil
}
