package worker

import (
	"context"
	"log"
	"path/filepath"

	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoProcessor struct {
	repoService     *service.RepoService
	githubService   *service.GitHubService
	fileTreeService *service.FileTreeService
	storagePath     string
	db              *pgxpool.Pool
}

func NewRepoProcessor(repoService *service.RepoService, githubService *service.GitHubService, fileTreeService *service.FileTreeService, storagePath string, db *pgxpool.Pool) *RepoProcessor {
	return &RepoProcessor{
		repoService:     repoService,
		githubService:   githubService,
		fileTreeService: fileTreeService,
		storagePath:     storagePath,
		db:              db,
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

	// Index the file tree into the database
	if p.fileTreeService != nil {
		if err := p.fileTreeService.IndexRepo(ctx, repoID, clonePath); err != nil {
			log.Printf("Failed to index file tree: %v", err)
		} else {
			log.Printf("Successfully indexed file tree for repo: %s", repoID)
		}
	}

	if err := p.parseDependencies(ctx, repoID, clonePath); err != nil {
		log.Printf("Failed to parse dependencies: %v", err)
	}

	//update status to completed (for now, parsing later....)
	if err := p.repoService.UpdateRepoStatus(ctx, repoID, models.StatusCompleted); err != nil {
		log.Printf("failed to update repo status: %v", err)
		return
	}

}

func (p *RepoProcessor) parseDependencies(ctx context.Context, repoID, repoPath string) error {
	parser := service.NewDependencyParser()
	deps, err := parser.ParseRepo(repoPath, repoID)
	if err != nil {
		return err
	}

	if len(deps) > 0 {
		depRepo := repository.NewDependencyRepository(p.db)
		if err := depRepo.CreateBatch(ctx, deps); err != nil {
			return err
		}
		log.Printf("Parsed %d dependencies for repo %s", len(deps), repoID)
	}
	return nil
}

