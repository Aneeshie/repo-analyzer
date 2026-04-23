package repository

import (
	"context"
	"os"
	"os/exec"
	"testing"

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
