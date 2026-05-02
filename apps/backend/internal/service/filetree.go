package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

// Directories to skip when indexing a repo's file tree.
var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	".next":        true,
	"__pycache__":  true,
	".turbo":       true,
	"dist":         true,
	"build":        true,
	"vendor":       true,
	"target":       true, // Rust/Java build output
}

// Extension-to-language mapping for syntax highlighting.
var extToLanguage = map[string]string{
	".go":         "go",
	".ts":         "typescript",
	".tsx":        "tsx",
	".js":         "javascript",
	".jsx":        "jsx",
	".py":         "python",
	".rs":         "rust",
	".java":       "java",
	".rb":         "ruby",
	".php":        "php",
	".c":          "c",
	".cpp":        "cpp",
	".h":          "c",
	".cs":         "csharp",
	".swift":      "swift",
	".kt":         "kotlin",
	".scala":      "scala",
	".html":       "html",
	".css":        "css",
	".scss":       "scss",
	".json":       "json",
	".yaml":       "yaml",
	".yml":        "yaml",
	".toml":       "toml",
	".xml":        "xml",
	".md":         "markdown",
	".sql":        "sql",
	".sh":         "bash",
	".bash":       "bash",
	".zsh":        "bash",
	".dockerfile": "dockerfile",
	".proto":      "protobuf",
	".graphql":    "graphql",
	".env":        "plaintext",
	".txt":        "plaintext",
	".csv":        "plaintext",
	".lock":       "plaintext",
	".mod":        "go",
	".sum":        "plaintext",
	".gitignore":  "plaintext",
	".editorconfig": "plaintext",
}

// MaxFileSize is the largest file (in bytes) we'll serve content for.
const MaxFileSize = 1 * 1024 * 1024 // 1MB

type FileTreeService struct {
	repo *repository.FileTreeRepository
}

func NewFileTreeService(repo *repository.FileTreeRepository) *FileTreeService {
	return &FileTreeService{repo: repo}
}

// IndexRepo walks a cloned repo directory and inserts every file/directory
// into the repo_files table.
func (s *FileTreeService) IndexRepo(ctx context.Context, repoID, repoPath string) error {
	var nodes []models.FileNode

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable files
		}

		// Get the path relative to the repo root
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return nil
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Skip ignored directories
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}

		// Compute parent path
		var parentPath *string
		dir := filepath.Dir(relPath)
		if dir != "." {
			parentPath = &dir
		}

		nodeType := "file"
		if info.IsDir() {
			nodeType = "directory"
		}

		// Detect language from extension
		var language *string
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if lang, ok := extToLanguage[ext]; ok {
				language = &lang
			}
			// Handle special filenames without extensions
			if language == nil {
				switch strings.ToLower(info.Name()) {
				case "dockerfile", "containerfile":
					l := "dockerfile"
					language = &l
				case "makefile":
					l := "makefile"
					language = &l
				}
			}
		}

		nodes = append(nodes, models.FileNode{
			RepoID:     repoID,
			ParentPath: parentPath,
			Path:       relPath,
			Name:       info.Name(),
			Type:       nodeType,
			Size:       info.Size(),
			Language:   language,
		})

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk repo directory: %w", err)
	}

	return s.repo.CreateBatch(ctx, nodes)
}

// GetTree returns all file nodes for a repo and assembles them into a nested tree.
func (s *FileTreeService) GetTree(ctx context.Context, repoID string) ([]*models.FileTreeNode, error) {
	nodes, err := s.repo.GetAll(ctx, repoID)
	if err != nil {
		return nil, err
	}

	return buildTree(nodes), nil
}

// GetChildren returns the direct children of a directory for lazy-loading.
func (s *FileTreeService) GetChildren(ctx context.Context, repoID, parentPath string) ([]models.FileNode, error) {
	return s.repo.GetChildren(ctx, repoID, parentPath)
}

// ReadFileContent reads a file's content from disk.
func (s *FileTreeService) ReadFileContent(repoPath, filePath string) (*models.FileContent, error) {
	fullPath := filepath.Join(repoPath, filePath)

	// Security: prevent path traversal attacks
	absRepo, _ := filepath.Abs(repoPath)
	absFile, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFile, absRepo) {
		return nil, fmt.Errorf("invalid file path")
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file")
	}

	if info.Size() > MaxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", info.Size(), MaxFileSize)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detect language
	var language *string
	ext := strings.ToLower(filepath.Ext(info.Name()))
	if lang, ok := extToLanguage[ext]; ok {
		language = &lang
	}

	return &models.FileContent{
		Path:     filePath,
		Name:     info.Name(),
		Content:  string(content),
		Language: language,
		Size:     info.Size(),
	}, nil
}

// buildTree assembles flat FileNode rows into a nested tree structure.
func buildTree(nodes []models.FileNode) []*models.FileTreeNode {
	// Map path -> tree node
	nodeMap := make(map[string]*models.FileTreeNode)

	// First pass: create all tree nodes
	for _, n := range nodes {
		nodeMap[n.Path] = &models.FileTreeNode{
			Name:     n.Name,
			Path:     n.Path,
			Type:     n.Type,
			Size:     n.Size,
			Language: n.Language,
		}
	}

	// Second pass: wire parent-child relationships
	var roots []*models.FileTreeNode
	for _, n := range nodes {
		treeNode := nodeMap[n.Path]
		if n.ParentPath == nil {
			roots = append(roots, treeNode)
		} else if parent, ok := nodeMap[*n.ParentPath]; ok {
			parent.Children = append(parent.Children, treeNode)
		} else {
			// Orphan node — add to roots as fallback
			roots = append(roots, treeNode)
		}
	}

	return roots
}
