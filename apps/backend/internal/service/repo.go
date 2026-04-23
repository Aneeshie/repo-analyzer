package service

import (
	"context"

	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type RepoService struct {
	repoRepo *repository.RepoRepository
}

func NewRepoService(repoRepo *repository.RepoRepository) *RepoService {
	return &RepoService{
		repoRepo: repoRepo,
	}
}

func (s *RepoService) CreateRepo(ctx context.Context, url string) (*models.Repo, error) {
	// validate URL format, check for duplicates, etc.

	return s.repoRepo.Create(ctx, url)
}

// GetRepo retrieves a repo by its ID.
func (s *RepoService) GetRepo(ctx context.Context, id string) (*models.Repo, error) {
	return s.repoRepo.FindByID(ctx, id)
}

// UpdateRepoStatus updates the status of a repo by its ID.
func (s *RepoService) UpdateRepoStatus(ctx context.Context, id string, status string) (*models.Repo, error) {
	return s.repoRepo.UpdateStatus(ctx, id, status)
}
