package repository

import (
	"context"
	"testing"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependencyRepository_CreateBatch(t *testing.T) {
	pool := setupTestDB(t)
	depRepo := NewDependencyRepository(pool)
	repoRepo := NewRepoRepository(pool)

	// Create a repo first
	repo, err := repoRepo.Create(context.Background(), "https://github.com/test/repo")
	require.NoError(t, err)

	deps := []models.Dependency{
		{
			RepoID:     repo.ID,
			Name:       "express",
			Version:    "^4.18.0",
			Ecosystem:  models.EcosystemNPM,
			Scope:      models.ScopeProduction,
			SourceFile: "package.json",
		},
		{
			RepoID:     repo.ID,
			Name:       "jest",
			Version:    "^29.0.0",
			Ecosystem:  models.EcosystemNPM,
			Scope:      models.ScopeDevelopment,
			SourceFile: "package.json",
		},
	}

	err = depRepo.CreateBatch(context.Background(), deps)
	assert.NoError(t, err)

	// Verify they were saved
	saved, err := depRepo.GetByRepoID(context.Background(), repo.ID)
	assert.NoError(t, err)
	assert.Len(t, saved, 2)
}
func TestDependencyRepository_CreateBatch_Empty(t *testing.T) {
	pool := setupTestDB(t)
	depRepo := NewDependencyRepository(pool)

	// Empty batch should not error
	err := depRepo.CreateBatch(context.Background(), []models.Dependency{})
	assert.NoError(t, err)
}

func TestDependencyRepository_CreateBatch_Duplicate(t *testing.T) {
	pool := setupTestDB(t)
	depRepo := NewDependencyRepository(pool)
	repoRepo := NewRepoRepository(pool)

	repo, err := repoRepo.Create(context.Background(), "https://github.com/test/repo")
	require.NoError(t, err)

	dep := models.Dependency{
		RepoID:     repo.ID,
		Name:       "express",
		Version:    "^4.18.0",
		Ecosystem:  models.EcosystemNPM,
		Scope:      models.ScopeProduction,
		SourceFile: "package.json",
	}

	// Insert same dep twice
	err = depRepo.CreateBatch(context.Background(), []models.Dependency{dep})
	assert.NoError(t, err)

	err = depRepo.CreateBatch(context.Background(), []models.Dependency{dep})
	assert.NoError(t, err) // Should not error (ON CONFLICT DO NOTHING)

	// Should only have one
	saved, err := depRepo.GetByRepoID(context.Background(), repo.ID)
	assert.NoError(t, err)
	assert.Len(t, saved, 1)
}
