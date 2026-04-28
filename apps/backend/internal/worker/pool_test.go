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

func setupTestPool(t *testing.T) (*Pool, *repository.RepoRepository, *pgxpool.Pool, func()) {
	// Setup test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	// Clean tables
	_, err = pool.Exec(context.Background(), "TRUNCATE repos, dependencies CASCADE")
	require.NoError(t, err)

	repoRepo := repository.NewRepoRepository(pool)
	repoService := service.NewRepoService(repoRepo)
	githubService := service.NewGitHubService()

	storagePath, err := os.MkdirTemp("", "worker-test-*")
	require.NoError(t, err)

	// Pass the db pool to NewPool
	workerPool := NewPool(repoService, githubService, storagePath, pool, 2)

	cleanup := func() {
		// Clean up after test
		_, _ = pool.Exec(context.Background(), "TRUNCATE repos, dependencies CASCADE")
		pool.Close()
		os.RemoveAll(storagePath)
	}

	return workerPool, repoRepo, pool, cleanup
}

func TestPool_AddJob(t *testing.T) {
	// Skip in CI due to timing issues
	if os.Getenv("CI") != "" {
		t.Skip("Skipping flaky test in CI")
	}

	workerPool, repoRepo, _, cleanup := setupTestPool(t)
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
	time.Sleep(10 * time.Second)

	// Check status updated
	updated, err := repoRepo.FindByID(context.Background(), repo.ID)
	require.NoError(t, err)

	// Should be completed (not pending or cloning)
	assert.NotEqual(t, models.StatusPending, updated.Status)
	assert.NotEqual(t, models.StatusCloning, updated.Status)
}

func TestPool_Shutdown(t *testing.T) {
	workerPool, _, _, cleanup := setupTestPool(t)
	defer cleanup()

	// Shutdown should not panic
	workerPool.Shutdown()
}

func TestPool_MultipleWorkers(t *testing.T) {
	workerPool, repoRepo, pool, cleanup := setupTestPool(t)
	defer cleanup()

	// Create multiple repos
	urls := []string{
		"https://github.com/octocat/Hello-World",
		"https://github.com/octocat/Spoon-Knife",
		"https://github.com/octocat/linguist",
	}

	repoIDs := make([]string, 0, len(urls))

	for _, url := range urls {
		repo, err := repoRepo.Create(context.Background(), url)
		require.NoError(t, err)
		repoIDs = append(repoIDs, repo.ID)

		workerPool.AddJob(models.Job{
			RepoID:  repo.ID,
			RepoURL: repo.URL,
		})
	}

	// Let workers process (increased timeout for CI)
	time.Sleep(30 * time.Second)

	// Verify all repos are processed
	for _, repoID := range repoIDs {
		var status string
		err := pool.QueryRow(context.Background(),
			"SELECT status FROM repos WHERE id = $1", repoID).Scan(&status)
		require.NoError(t, err)
		assert.NotEqual(t, models.StatusPending, status)
		assert.NotEqual(t, models.StatusCloning, status)
	}
}
