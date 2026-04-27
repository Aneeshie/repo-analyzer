package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	client *github.Client
}

func NewGitHubService() *GitHubService {
	var client *github.Client

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &GitHubService{
		client: client,
	}
}

// ParseGitHubURL extracts owner and repo name from URL
func (s *GitHubService) ParseGitHubURL(url string) (owner, repo string, err error) {
	// Remove https:// or http://
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Remove trailing slash if present
	url = strings.TrimSuffix(url, "/")

	// Split by slash
	parts := strings.Split(url, "/")

	// Need at least 2 parts: domain + owner + repo
	// Example: github.com/facebook/react -> 3 parts
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid GitHub URL: need owner and repo name")
	}

	// Last two parts should be owner/repo
	// Skip the domain (parts[0])
	owner = parts[len(parts)-2]
	repo = strings.TrimSuffix(parts[len(parts)-1], ".git")

	// Validate owner and repo are not empty
	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("invalid GitHub URL: owner or repo is empty")
	}

	return owner, repo, nil
}

// CloneRepo clones a GitHub repository to local path
func (s *GitHubService) CloneRepo(ctx context.Context, repoURL, localPath string) error {
	// Shallow clone (depth 1) for speed
	_, err := git.PlainCloneContext(ctx, localPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: nil,
		Depth:    1,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	return nil
}

// GetRepoMetadata fetches repo info from GitHub API
func (s *GitHubService) GetRepoMetadata(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo metadata: %w", err)
	}
	return repository, nil
}
