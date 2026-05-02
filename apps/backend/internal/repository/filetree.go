package repository

import (
	"context"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileTreeRepository struct {
	db *pgxpool.Pool
}

func NewFileTreeRepository(db *pgxpool.Pool) *FileTreeRepository {
	return &FileTreeRepository{db: db}
}

// CreateBatch inserts all file nodes for a repo in one go.
func (r *FileTreeRepository) CreateBatch(ctx context.Context, files []models.FileNode) error {
	if len(files) == 0 {
		return nil
	}

	query := `
	INSERT INTO repo_files (repo_id, parent_path, path, name, type, size, language)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (repo_id, path) DO NOTHING
	`

	for _, f := range files {
		_, err := r.db.Exec(ctx, query,
			f.RepoID, f.ParentPath, f.Path, f.Name, f.Type, f.Size, f.Language)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetChildren returns the direct children of a given directory path.
// Pass parentPath="" to get root-level entries.
func (r *FileTreeRepository) GetChildren(ctx context.Context, repoID string, parentPath string) ([]models.FileNode, error) {
	var query string
	//var rows interface{ Next() bool }
	var err error

	if parentPath == "" {
		query = `
			SELECT id, repo_id, parent_path, path, name, type, size, language
			FROM repo_files
			WHERE repo_id = $1 AND parent_path IS NULL
			ORDER BY type DESC, name ASC
		`
		rows2, err2 := r.db.Query(ctx, query, repoID)
		if err2 != nil {
			return nil, err2
		}
		defer rows2.Close()
		return scanFileNodes(rows2)
	}

	query = `
		SELECT id, repo_id, parent_path, path, name, type, size, language
		FROM repo_files
		WHERE repo_id = $1 AND parent_path = $2
		ORDER BY type DESC, name ASC
	`
	pgxRows, err := r.db.Query(ctx, query, repoID, parentPath)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()
	return scanFileNodes(pgxRows)
}

// GetAll returns every file node for a repo (used to build the full tree).
func (r *FileTreeRepository) GetAll(ctx context.Context, repoID string) ([]models.FileNode, error) {
	query := `
		SELECT id, repo_id, parent_path, path, name, type, size, language
		FROM repo_files
		WHERE repo_id = $1
		ORDER BY path ASC
	`

	rows, err := r.db.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFileNodes(rows)
}

type scannable interface {
	Next() bool
	Scan(dest ...any) error
}

func scanFileNodes(rows scannable) ([]models.FileNode, error) {
	var nodes []models.FileNode
	for rows.Next() {
		var node models.FileNode
		err := rows.Scan(
			&node.ID, &node.RepoID, &node.ParentPath,
			&node.Path, &node.Name, &node.Type,
			&node.Size, &node.Language,
		)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
