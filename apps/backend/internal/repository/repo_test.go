package repository

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer_test?sslmode=disable"
	}

	// Get project root using git
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}
	projectRoot := string(output)
	projectRoot = projectRoot[:len(projectRoot)-1] // Remove newline

	// Run migrations from project root
	migrateCmd := exec.Command("migrate",
		"-path", projectRoot+"/infra/migrations",
		"-database", dbURL,
		"up")
	migrateCmd.Dir = projectRoot

	if err := migrateCmd.Run(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		// Truncate all tables
		_, _ = pool.Exec(context.Background(), `
			TRUNCATE repos, repo_details, analysis CASCADE;
		`)
		pool.Close()
	})

	return pool
}

func TestRepoRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepoRepository(pool)

	url := "https://github.com/test/repo"
	result, err := repo.Create(context.Background(), url)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, url, result.URL)
	assert.Equal(t, "pending", result.Status)
}

func TestRepoRepository_CreateDuplicate(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepoRepository(pool)
	url := "https://github.com/test/duplicate"

	_, err := repo.Create(context.Background(), url)
	assert.NoError(t, err)

	_, err = repo.Create(context.Background(), url)
	assert.Error(t, err)
}

func TestRepoRepository_FindByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepoRepository(pool)

	// First create a repo
	url := "https://github.com/test/findme"
	created, err := repo.Create(context.Background(), url)
	require.NoError(t, err)

	// Then find it
	found, err := repo.FindByID(context.Background(), created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, url, found.URL)
	assert.Equal(t, models.StatusPending, found.Status)
}

func TestRepoRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepoRepository(pool)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.Error(t, err)
}

func TestRepoRepository_UpdateStatus(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepoRepository(pool)

	// Create a repo
	url := "https://github.com/test/updateme"
	created, err := repo.Create(context.Background(), url)
	require.NoError(t, err)

	// Update status
	err = repo.UpdateStatus(context.Background(), created.ID, models.StatusCloning)
	require.NoError(t, err)

	// Verify status changed
	updated, err := repo.FindByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusCloning, updated.Status)
}

func TestRepoRepository_UpdateStatus_InvalidID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepoRepository(pool)

	err := repo.UpdateStatus(context.Background(), "00000000-0000-0000-0000-000000000000", models.StatusCompleted)

	assert.Error(t, err)
}
