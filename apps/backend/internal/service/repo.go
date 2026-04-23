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
