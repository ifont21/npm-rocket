package adapters

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v52/github"
	"github.com/ifont21/pre-releaser-cli/internal/domain"
	"golang.org/x/oauth2"
)

func getAffectedPackagesForDescription(affectedPackages []string) string {
	affectedPackagesDescription := ""
	for _, affectedPackage := range affectedPackages {
		affectedPackagesDescription += fmt.Sprintf("- %s\n", affectedPackage)
	}
	return affectedPackagesDescription
}

type PRGithubRepository struct {
	config Config
}

func NewPRGithubRepository(config Config) *PRGithubRepository {
	return &PRGithubRepository{
		config: config,
	}
}

func (p *PRGithubRepository) CreatePR(pr domain.PR, affectedPackages []string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repoConfig, err := p.config.GetRepositoryConfig()
	if err != nil {
		return err
	}

	repository := repoConfig.Name
	owner := repoConfig.Owner
	baseBranch := repoConfig.Branch
	prTitleDesc := fmt.Sprintf("## Prepare release for NPM packages in **%s** repository", repository)
	prAffectedHeader := "### Affected packages - version and CHANGELOG updated"
	prAffectedPackages := getAffectedPackagesForDescription(affectedPackages)
	prDescription := fmt.Sprintf("%s\n%s\n%s", prTitleDesc, prAffectedHeader, prAffectedPackages)

	_, _, err = client.PullRequests.Create(ctx, owner, repository, &github.NewPullRequest{
		Title: &pr.TitlePR,
		Head:  &pr.BranchName,
		Base:  &baseBranch,
		Body:  &prDescription,
	})
	if err != nil {
		fmt.Println("Error creating PR", err)
		return err
	}
	fmt.Println("PR created ************")

	return nil
}
