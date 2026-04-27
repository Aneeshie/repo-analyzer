package service

import (
	"context"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type RepoServiceInterface interface {
	CreateRepo(ctx context.Context, url string) (*models.Repo, error)
	GetRepo(ctx context.Context, id string) (*models.Repo, error)
	UpdateRepoStatus(ctx context.Context, id string, status string) error
}
