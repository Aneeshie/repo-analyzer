package repository

import (
	"context"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DependencyRepository struct {
	db *pgxpool.Pool
}

func NewDependencyRepository(db *pgxpool.Pool) *DependencyRepository {
	return &DependencyRepository{db}
}

func (r *DependencyRepository) CreateBatch(ctx context.Context, deps []models.Dependency) error {
	if len(deps) == 0 {
		return nil
	}

	query := `
	INSERT INTO dependencies (repo_id, name, version, ecosystem, scope, source_file)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (repo_id, name, ecosystem, source_file) DO NOTHING
	`

	for _, dep := range deps {
		_, err := r.db.Exec(ctx, query,
			dep.RepoID, dep.Name, dep.Version, dep.Ecosystem, dep.Scope, dep.SourceFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DependencyRepository) GetByRepoID(ctx context.Context, repoID string) ([]models.Dependency, error) {
	query := `
		SELECT id, repo_id, name, version, ecosystem, scope, source_file, created_at
		FROM dependencies
		WHERE repo_id = $1
		ORDER BY ecosystem, name
	`

	rows, err := r.db.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []models.Dependency
	for rows.Next() {
		var dep models.Dependency
		err := rows.Scan(
			&dep.ID, &dep.RepoID, &dep.Name, &dep.Version,
			&dep.Ecosystem, &dep.Scope, &dep.SourceFile, &dep.CreatedAt)
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, nil
}
