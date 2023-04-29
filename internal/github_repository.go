package internal

type GithubRepository interface {
	ListRepositories() ([]string, error)
	CreatePR(filePath string) error
}
