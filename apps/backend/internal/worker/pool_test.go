package worker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDBLockID int64 = 928374650192 // stable lock id to serialize DB-mutating tests

func acquireTestDBLock(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), "SELECT pg_advisory_lock($1)", testDBLockID)
	require.NoError(t, err)
}

func releaseTestDBLock(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, _ = pool.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", testDBLockID)
}

func waitForRepoProcessed(t *testing.T, repoRepo *repository.RepoRepository, repoID string, timeout time.Duration) string {
	t.Helper()
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		repo, err := repoRepo.FindByID(context.Background(), repoID)
		if err == nil {
			if repo.Status != models.StatusPending && repo.Status != models.StatusCloning {
				return repo.Status
			}
		} else if !errors.Is(err, pgx.ErrNoRows) {
			require.NoError(t, err)
		}

		time.Sleep(250 * time.Millisecond)
	}

	repo, err := repoRepo.FindByID(context.Background(), repoID)
	if err != nil {
		require.NoError(t, fmt.Errorf("timed out waiting for repo %s to be processed: %w", repoID, err))
	}
	return repo.Status
}

func setupTestPool(t *testing.T) (*Pool, *repository.RepoRepository, *pgxpool.Pool, func()) {
	// Setup test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	// Serialize tests that mutate shared tables so parallel package tests
	// don't TRUNCATE each other's data.
	acquireTestDBLock(t, pool)

	// Clean tables
	_, err = pool.Exec(context.Background(), "TRUNCATE repos, dependencies CASCADE")
	require.NoError(t, err)

	repoRepo := repository.NewRepoRepository(pool)
	depRepo := repository.NewDependencyRepository(pool)
	repoService := service.NewRepoService(repoRepo, depRepo)
	githubService := service.NewGitHubService()

	storagePath, err := os.MkdirTemp("", "worker-test-*")
	require.NoError(t, err)

	// Pass the db pool to NewPool
	workerPool := NewPool(repoService, githubService, nil, storagePath, pool, 2)

	cleanup := func() {
		// Clean up after test
		_, _ = pool.Exec(context.Background(), "TRUNCATE repos, dependencies CASCADE")
		releaseTestDBLock(t, pool)
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

	// Check status updated
	status := waitForRepoProcessed(t, repoRepo, repo.ID, 20*time.Second)

	// Should be completed (not pending or cloning)
	assert.NotEqual(t, models.StatusPending, status)
	assert.NotEqual(t, models.StatusCloning, status)
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

	// Verify all repos are processed
	for _, repoID := range repoIDs {
		status := waitForRepoProcessed(t, repoRepo, repoID, 60*time.Second)
		assert.NotEqual(t, models.StatusPending, status)
		assert.NotEqual(t, models.StatusCloning, status)
	}

	// Ensure pool is still usable (sanity check that connection isn't wedged).
	require.NoError(t, pool.Ping(context.Background()))
}

