package domain

type PR struct {
	TitlePR     string
	BranchName  string
	Base        string
	Description string
}

type PRRepository interface {
	CreatePR(pullRequest PR, affectedPackages []string) error
}
