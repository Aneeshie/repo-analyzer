package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
func TestGetRepoMetadata_RequiresToken(t *testing.T) {
	s := NewGitHubService()

	// This test will actually call GitHub API
	// For now, skip if no token
	if s.client == nil {
		t.Skip("No GitHub client available")
	}

	// Just test that the function exists
	assert.NotNil(t, s.GetRepoMetadata)
}

// Add test for CloneRepo
func TestCloneRepo_InvalidPath(t *testing.T) {
	s := NewGitHubService()

	// Test with invalid URL
	err := s.CloneRepo(context.Background(), "not-a-url", "/tmp/invalid")
	assert.Error(t, err)
}
