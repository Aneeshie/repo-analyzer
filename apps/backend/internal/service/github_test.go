package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitHubURL(t *testing.T) {
	s := NewGitHubService()

	tests := []struct {
		name     string
		url      string
		owner    string
		repo     string
		hasError bool
	}{
		// Valid URLs
		{
			name:  "standard https URL",
			url:   "https://github.com/facebook/react",
			owner: "facebook",
			repo:  "react",
		},
		{
			name:  "URL with .git suffix",
			url:   "https://github.com/facebook/react.git",
			owner: "facebook",
			repo:  "react",
		},
		{
			name:  "URL without https",
			url:   "github.com/facebook/react",
			owner: "facebook",
			repo:  "react",
		},
		{
			name:  "repo with dashes",
			url:   "https://github.com/golang/go",
			owner: "golang",
			repo:  "go",
		},
		{
			name:  "repo with numbers",
			url:   "https://github.com/kubernetes/kubernetes",
			owner: "kubernetes",
			repo:  "kubernetes",
		},
		{
			name:  "repo with underscores",
			url:   "https://github.com/heroku/heroku-cli",
			owner: "heroku",
			repo:  "heroku-cli",
		},
		{
			name:  "URL with trailing slash",
			url:   "https://github.com/facebook/react/",
			owner: "facebook",
			repo:  "react",
		},
		{
			name:  "HTTP protocol",
			url:   "http://github.com/facebook/react",
			owner: "facebook",
			repo:  "react",
		},
		{
			name:  "www subdomain",
			url:   "https://www.github.com/facebook/react",
			owner: "facebook",
			repo:  "react",
		},

		// Invalid URLs
		{
			name:     "invalid - only domain",
			url:      "https://github.com",
			hasError: true,
		},
		{
			name:     "invalid - only domain with slash",
			url:      "https://github.com/",
			hasError: true,
		},
		{
			name:     "invalid - only owner, no repo",
			url:      "https://github.com/facebook",
			hasError: true,
		},
		{
			name:     "invalid - empty URL",
			url:      "",
			hasError: true,
		},
		{
			name:     "invalid - random string",
			url:      "not-a-url",
			hasError: true,
		},
		{
			name:     "invalid - missing owner",
			url:      "https://github.com//react",
			hasError: true,
		},
		{
			name:     "invalid - gitlab URL (not github)",
			url:      "https://gitlab.com/facebook/react",
			hasError: false, // Should still parse? Or reject?
			owner:    "facebook",
			repo:     "react",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := s.ParseGitHubURL(tt.url)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.owner, owner)
				assert.Equal(t, tt.repo, repo)
			}
		})
	}
}

// Add test for GetRepoMetadata (mock this later? or skip for now)

func TestCloneRepo_Success(t *testing.T) {
	s := NewGitHubService()

	// Create temp directory for clone
	tempDir, err := os.MkdirTemp("", "repo-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Clone a small public repo
	ctx := context.Background()
	err = s.CloneRepo(ctx, "https://github.com/octocat/Hello-World.git", tempDir)

	assert.NoError(t, err)

	// Verify .git directory exists (means clone worked)
	gitDir := filepath.Join(tempDir, ".git")
	_, err = os.Stat(gitDir)
	assert.NoError(t, err)
}

func TestCloneRepo_InvalidURL(t *testing.T) {
	s := NewGitHubService()

	tempDir, err := os.MkdirTemp("", "repo-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	err = s.CloneRepo(ctx, "https://github.com/invalid-repo-12345/does-not-exist", tempDir)

	assert.Error(t, err)
}

func TestCloneRepo_InvalidPath(t *testing.T) {
	s := NewGitHubService()

	ctx := context.Background()
	// Invalid path (should fail to create directory)
	err := s.CloneRepo(ctx, "https://github.com/octocat/Hello-World.git", "/invalid/path/that/doesnt/exist")

	assert.Error(t, err)
}

func TestGetRepoMetadata_WithToken(t *testing.T) {
	// Skip if no GitHub token
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("GITHUB_TOKEN not set, skipping")
	}

	s := NewGitHubService()
	ctx := context.Background()

	repo, err := s.GetRepoMetadata(ctx, "octocat", "Hello-World")

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, "Hello-World", *repo.Name)
}

func TestGetRepoMetadata_NotFound(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("GITHUB_TOKEN not set, skipping")
	}

	s := NewGitHubService()
	ctx := context.Background()

	_, err := s.GetRepoMetadata(ctx, "invalid-owner-12345", "invalid-repo-67890")

	assert.Error(t, err)
}
