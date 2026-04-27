package service

import (
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
			name:     "invalid URL - too short",
			url:      "https://github.com/facebook",
			hasError: true,
		},
		{
			name:     "empty URL",
			url:      "",
			hasError: true,
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
