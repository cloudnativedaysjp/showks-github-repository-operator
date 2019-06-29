package gh

import (
	"context"
	"github.com/google/go-github/github"
)

type GitHubClientInterface interface {
	CreateRepository(org string, repo *github.Repository) (*github.Repository, error)
	GetRepository(org string, repo string) (*github.Repository, error)
	AddCollaborator(owner string, repo string, user string, permission string) error
	GetPermissionLevel(owner string, repo string, user string) (string, error)
	UpdateBranchProtection(owner string, repo string, branch string, request *github.ProtectionRequest) (*github.Protection, error)
	CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, error)
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
		return nil, &NotFoundError{}
	}

	return repo, err
}

func (c *GithubClient) AddCollaborator(owner string, repo string, user string, permission string) error {
	opt := &github.RepositoryAddCollaboratorOptions{Permission: permission}
	ctx := context.Background()
	_, err := c.client.Repositories.AddCollaborator(ctx, owner, repo, user, opt)
	return err
}

func (c *GithubClient) GetPermissionLevel(owner string, repo string, user string) (string, error) {
	ctx := context.Background()
	pl, resp, err := c.client.Repositories.GetPermissionLevel(ctx, owner, repo, user)
	if err != nil {
		return "", err
	}
	if resp.Response.StatusCode == 404 {
		return "", &NotFoundError{}
	}
	return *pl.Permission, nil
}

func (c *GithubClient) UpdateBranchProtection(owner string, repo string, branch string, request *github.ProtectionRequest) (*github.Protection, error) {
	ctx := context.Background()
	p, _, err := c.client.Repositories.UpdateBranchProtection(ctx, owner, repo, branch, request)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (c *GithubClient) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, error) {
	ctx := context.Background()
	h, _, err := c.client.Repositories.CreateHook(ctx, owner, repo, hook)
	if err != nil {
		return h, err
	}

	return h, nil
}

type NotFoundError struct {
	error
}

func (e *NotFoundError) Error() string {
	return e.error.Error()
}
