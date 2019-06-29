package gh

import (
	"context"
	"github.com/google/go-github/github"
)

type GitHubClientInterface interface {
	CreateRepository(org string, repo *github.Repository) (*github.Repository, error)
	GetRepository(org string, repo string) (*github.Repository, error)
	AddCoraborator(owner string, repo string, user string, permission string) error
}

func NewClient() GitHubClientInterface {
	c := github.NewClient(nil)
	return &GithubClient{client: c}
}

type GithubClient struct {
	client *github.Client
}

func (c *GithubClient) CreateRepository(org string, repo *github.Repository) (*github.Repository, error) {
	ctx := context.Background()
	repo, _, err := c.client.Repositories.Create(ctx, org, repo)

	return repo, err
}

func (c *GithubClient) GetRepository(org string, repoName string) (*github.Repository, error) {
	ctx := context.Background()
	repo, resp, err := c.client.Repositories.Get(ctx, org, repoName)

	if resp.Response.StatusCode == 404 {
		return nil, NotFoundError{}
	}

	return repo, err
}

func (c *GithubClient) AddCoraborator(owner string, repo string, user string, permission string) error {
	opt := &github.RepositoryAddCollaboratorOptions{Permission: permission}
	ctx := context.Background()
	_, err := c.client.Repositories.AddCollaborator(ctx, owner, repo, user, opt)
	return err
}

type NotFoundError struct {
	error
}
