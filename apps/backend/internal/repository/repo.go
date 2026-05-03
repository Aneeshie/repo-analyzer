package repository

import (
	"context"
	"fmt"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoRepository struct {
	db *pgxpool.Pool
}

func NewRepoRepository(db *pgxpool.Pool) *RepoRepository {
	return &RepoRepository{
		db: db,
	}
}

func (r *RepoRepository) Create(ctx context.Context, url string) (*models.Repo, error) {
	var repo models.Repo

	query := `INSERT INTO repos (url) VALUES ($1)
			  RETURNING id, url, status, created_at, updated_at, entry_points
	`

	err := r.db.QueryRow(ctx, query, url).Scan(&repo.ID, &repo.URL, &repo.Status, &repo.CreatedAt, &repo.UpdatedAt, &repo.EntryPoints)

	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *RepoRepository) FindByID(ctx context.Context, id string) (*models.Repo, error) {
	var repo models.Repo

	query := `SELECT id, url, status, created_at, updated_at, entry_points FROM repos WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(&repo.ID, &repo.URL, &repo.Status, &repo.CreatedAt, &repo.UpdatedAt, &repo.EntryPoints)

	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *RepoRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE repos SET status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	// Check if any row was updated
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("repo not found")
	}

	return nil
}

func (r *RepoRepository) UpdateEntryPoints(ctx context.Context, id string, entryPoints []string) error {
	query := `UPDATE repos SET entry_points = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, entryPoints, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("repo not found")
	}

	return nil
}
