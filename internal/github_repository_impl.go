package internal

import (
	"context"
	"os"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

type GithubRepositoryImpl struct{}

func NewGithubRepositoryImpl() *GithubRepositoryImpl {
	return &GithubRepositoryImpl{}
}

func (g *GithubRepositoryImpl) ListRepositories() ([]string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListOptions{
		Type: "all",
	}
	repos, _, err := client.Repositories.List(context.Background(), "", opt)
	if err != nil {
		return nil, err
	}
	var repoNames []string
	for _, repo := range repos {
		repoNames = append(repoNames, *repo.Name)
	}
	return repoNames, nil
}

func (g *GithubRepositoryImpl) CreatePR(filePath string) error {
	/* ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	preReleaserYaml, err := files.GetReleaseCLIConfig(fmt.Sprintf("%s/pre-releaser.yaml", filePath))
	if err != nil {
		return err
	}

	owner := preReleaserYaml.Repository.Owner
	repository := preReleaserYaml.Repository.Name

	pr, _, err := client.PullRequests.Create(ctx, owner, repository, &github.NewPullRequest{
		Title: github.String("test: PR testing"),
		Head:  github.String("ci/releaser-cli-repo"),
		Base:  github.String("main"),
		Body:  github.String("## Description \n\n This is a test PR"),
	})
	if err != nil {
		fmt.Println("Error creating PR", err)
		return err
	}
	fmt.Println("PR created", pr) */

	return nil
}
