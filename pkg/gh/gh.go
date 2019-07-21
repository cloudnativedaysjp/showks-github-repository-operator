package gh

import (
	"fmt"
	"github.com/cloudnativedaysjp/showks-github-repository-operator/pkg/apis/showks/v1beta1"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io"
	"os"
	"time"
)

type GitHubClientInterface interface {
	CreateRepository(org string, repo *github.Repository) (*github.Repository, error)
	DeleteRepository(org string, repo string) error
	GetRepository(org string, repo string) (*github.Repository, error)
	InitializeRepository(rs v1beta1.GitHubRepositorySpec) error
	ListCollaborator(owner string, repo string) ([]*github.User, error)
	AddCollaborator(owner string, repo string, user string, permission string) error
	RemoveCollaborator(owner string, repo string, user string) error
	GetPermissionLevel(owner string, repo string, user string) (string, error)
	UpdateBranchProtection(owner string, repo string, branch string, request *github.ProtectionRequest) (*github.Protection, error)
	ListHook(owner string, repo string) ([]*github.Hook, error)
	CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, error)
	UpdateHook(owner string, repo string, id int64, hook *github.Hook) (*github.Hook, error)
	DeleteHook(owner string, repo string, id int64) error
	GetUser(user string) (*github.User, error)
}

func NewClient() GitHubClientInterface {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)
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

func (c *GithubClient) DeleteRepository(org string, repo string) error {
	ctx := context.Background()
	_, err := c.client.Repositories.Delete(ctx, org, repo)
	if err != nil {
		return err
	}

	return nil
}

func (c *GithubClient) GetRepository(org string, repoName string) (*github.Repository, error) {
	ctx := context.Background()
	repo, resp, err := c.client.Repositories.Get(ctx, org, repoName)

	if err != nil && resp != nil {
		if resp.Response.StatusCode == 404 {
			return nil, &NotFoundError{}
		}
	}

	if err != nil {
		return nil, err
	}

	return repo, err
}

func (c *GithubClient) InitializeRepository(rs v1beta1.GitHubRepositorySpec) error {
	org := rs.Org
	name := rs.Name
	f := memfs.New()
	githubToken := os.Getenv("GITHUB_TOKEN")
	user, err := c.GetUser("")
	if err != nil {
		return err
	}
	auth := &http.BasicAuth{Username: *user.Login, Password: githubToken}

	repo, err := git.Clone(memory.NewStorage(), f, &git.CloneOptions{
		URL:           "https://github.com/" + rs.RepositoryTemplate.Org + "/" + rs.RepositoryTemplate.Name + ".git",
		ReferenceName: plumbing.ReferenceName("refs/heads/master"),
		Auth:          auth,
	})
	if err != nil {
		return err
	}

	err = repo.DeleteRemote("origin")
	if err != nil {
		return err
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/" + org + "/" + name + ".git"},
	})
	if err != nil {
		return err
	}

	for _, ic := range rs.RepositoryTemplate.InitialCommits {
		a, err := f.Create(ic.Path)
		if err != nil {
			return errors.Wrapf(err, "Faield to f.Create")
		}
		defer a.Close()
		fmt.Printf("contents: %s\n", ic.Contents)
		_, err = io.WriteString(a, ic.Contents)
		if err != nil {
			return errors.Wrapf(err, "Faield to Copy")
		}
		a.Close()

		w, err := repo.Worktree()
		if err != nil {
			return errors.Wrapf(err, "Failed to Worktree")
		}

		_, err = w.Add(ic.Path)
		if err != nil {
			return errors.Wrapf(err, "Failed to Add")
		}

		_, err = w.Commit(fmt.Sprintf("Add %s", ic.Path), &git.CommitOptions{
			Author: &object.Signature{
				Name:  rs.RepositoryTemplate.Username,
				Email: rs.RepositoryTemplate.Email,
				When:  time.Now(),
			},
		})

		if err != nil {
			return errors.Wrapf(err, "Failed to Commit")
		}
	}

	for _, ib := range rs.RepositoryTemplate.InitialBranches {
		err = repo.Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs:   []config.RefSpec{config.RefSpec(ib)},
			Auth:       auth,
		})
		if err != nil {
			return errors.Wrapf(err, "Failed to Push")
		}
	}

	return nil
}

func (c *GithubClient) ListCollaborator(owner string, repo string) ([]*github.User, error) {
	opt := &github.ListCollaboratorsOptions{}
	ctx := context.Background()
	users, _, err := c.client.Repositories.ListCollaborators(ctx, owner, repo, opt)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *GithubClient) AddCollaborator(owner string, repo string, user string, permission string) error {
	opt := &github.RepositoryAddCollaboratorOptions{Permission: permission}
	ctx := context.Background()
	_, err := c.client.Repositories.AddCollaborator(ctx, owner, repo, user, opt)
	return err
}

func (c *GithubClient) RemoveCollaborator(owner string, repo string, user string) error {
	ctx := context.Background()
	_, err := c.client.Repositories.RemoveCollaborator(ctx, owner, repo, user)
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

func (c *GithubClient) ListHook(owner string, repo string) ([]*github.Hook, error) {
	ctx := context.Background()
	hooks, _, err := c.client.Repositories.ListHooks(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return hooks, err
	}

	return hooks, nil
}

func (c *GithubClient) UpdateHook(owner string, repo string, id int64, hook *github.Hook) (*github.Hook, error) {
	ctx := context.Background()
	hook, _, err := c.client.Repositories.EditHook(ctx, owner, repo, id, hook)
	if err != nil {
		return hook, err
	}

	return hook, nil
}

func (c *GithubClient) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, error) {
	ctx := context.Background()
	h, _, err := c.client.Repositories.CreateHook(ctx, owner, repo, hook)
	if err != nil {
		return h, err
	}

	return h, nil
}

func (c *GithubClient) DeleteHook(owner string, repo string, id int64) error {
	ctx := context.Background()
	_, err := c.client.Repositories.DeleteHook(ctx, owner, repo, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *GithubClient) GetUser(name string) (*github.User, error) {
	user, _, err := c.client.Users.Get(context.Background(), name)
	if err != nil {
		return user, err
	}
	return user, nil
}

type NotFoundError struct {
	error
}

func (e *NotFoundError) Error() string {
	return e.error.Error()
}
