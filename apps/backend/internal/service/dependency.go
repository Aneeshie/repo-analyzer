package service

import (
	"context"

	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type DependencyService struct{
	depRepo *repository.DependencyRepository

}

func NewDependencyService(depRepo *repository.DependencyRepository) *DependencyService {
	return &DependencyService{
		depRepo: depRepo,
	}
}

func (s *DependencyService) GetDependencyByID(ctx context.Context, repoId string) ([]models.Dependency, error){
	return s.depRepo.GetByRepoID(ctx, repoId)
}


