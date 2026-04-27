package models

import "time"

const (
	StatusPending   = "pending"
	StatusAnalyzing = "analyzing"
	StatusCompleted = "completed"
	StatusParsing   = "parsing"
	StatusFailed    = "failed"
	StatusCloning   = "cloning"
)

type Repo struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Job struct {
	RepoID  string
	RepoURL string
}

type CreateRepoRequest struct {
	URL string `json:"url"`
}

type CreateRepoResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

type GetRepoResponse struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
