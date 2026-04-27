package worker

import (
	"context"
	"log"
	"path/filepath"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type RepoProcessor struct {
	repoService   *service.RepoService
	githubService *service.GitHubService
	storagePath   string
}

func NewRepoProcessor(repoService *service.RepoService, githubService *service.GitHubService, storagePath string) *RepoProcessor {
	return &RepoProcessor{
		repoService:   repoService,
		githubService: githubService,
		storagePath:   storagePath,
	}
}

func (p *RepoProcessor) ProcessRepo(ctx context.Context, repoID, repoURL string) {
	if err := p.repoService.UpdateRepoStatus(ctx, repoID, models.StatusCloning); err != nil {
		log.Printf("failed to update repo status: %v", err)
		return
	}

	log.Printf("cloning repo: %s", repoURL)

	//clone repo

	clonePath := filepath.Join(p.storagePath, repoID)
	if err := p.githubService.CloneRepo(ctx, repoURL, clonePath); err != nil {
		log.Printf("failed to clone repo: %v", err)
		p.repoService.UpdateRepoStatus(ctx, repoID, models.StatusFailed)
		return
	}

	log.Printf("Successfully cloned repo: %s", repoURL)

	//update status to completed (for now, parsing later....)
	if err := p.repoService.UpdateRepoStatus(ctx, repoID, models.StatusCompleted); err != nil {
		log.Printf("failed to update repo status: %v", err)
		return
	}

}
