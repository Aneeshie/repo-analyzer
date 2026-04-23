package repository

import (
	"context"

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
			  RETURNING id, url, status, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, url).Scan(&repo.ID, &repo.URL, &repo.Status, &repo.CreatedAt, &repo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &repo, nil
}
