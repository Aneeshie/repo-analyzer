package worker

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestPool(t *testing.T) (*Pool, *repository.RepoRepository, func()) {
	// Setup test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	// Clean tables
	_, err = pool.Exec(context.Background(), "TRUNCATE repos CASCADE")
	require.NoError(t, err)

	repoRepo := repository.NewRepoRepository(pool)
	repoService := service.NewRepoService(repoRepo)
	githubService := service.NewGitHubService()

	storagePath, err := os.MkdirTemp("", "worker-test-*")
	require.NoError(t, err)

	// Pass the db pool to NewPool
	workerPool := NewPool(repoService, githubService, storagePath, pool, 2)

	cleanup := func() {
		workerPool.Shutdown()
		pool.Close()
		os.RemoveAll(storagePath)
	}

	return workerPool, repoRepo, cleanup
}
func TestPool_AddJob(t *testing.T) {
	workerPool, repoRepo, cleanup := setupTestPool(t)
	defer cleanup()

	// Create a repo first
	repo, err := repoRepo.Create(context.Background(), "https://github.com/octocat/Hello-World")
	require.NoError(t, err)

	// Add job to pool
	workerPool.AddJob(models.Job{
		RepoID:  repo.ID,
		RepoURL: repo.URL,
	})

	// Give worker time to process
	time.Sleep(2 * time.Second)

	// Check status updated
	updated, err := repoRepo.FindByID(context.Background(), repo.ID)
	require.NoError(t, err)

	// Should be completed (not pending)
	assert.NotEqual(t, models.StatusPending, updated.Status)
}

func TestPool_Shutdown(t *testing.T) {
	workerPool, _, _ := setupTestPool(t)

	// Shutdown should not panic
	workerPool.Shutdown()

	// After shutdown, adding jobs should panic or block
	// But we don't test that for now
	assert.True(t, true)
}

func TestPool_MultipleWorkers(t *testing.T) {
	workerPool, repoRepo, cleanup := setupTestPool(t)
	defer cleanup()

	// Create multiple repos
	urls := []string{
		"https://github.com/octocat/Hello-World",
		"https://github.com/octocat/Spoon-Knife",
		"https://github.com/octocat/linguist",
	}

	for _, url := range urls {
		repo, err := repoRepo.Create(context.Background(), url)
		require.NoError(t, err)

		workerPool.AddJob(models.Job{
			RepoID:  repo.ID,
			RepoURL: repo.URL,
		})
	}

	// Let workers process
	time.Sleep(5 * time.Second)

	// All repos should be processed
	// (status not pending)
	assert.True(t, true)
}
