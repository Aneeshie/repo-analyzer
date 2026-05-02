package models

// FileNode represents a single file or directory in a repo's file tree.
type FileNode struct {
	ID         string  `json:"id,omitempty"`
	RepoID     string  `json:"repo_id,omitempty"`
	ParentPath *string `json:"parent_path,omitempty"` // nil for root-level items
	Path       string  `json:"path"`
	Name       string  `json:"name"`
	Type       string  `json:"type"` // "file" or "directory"
	Size       int64   `json:"size,omitempty"`
	Language   *string `json:"language,omitempty"`
}

// FileTreeNode is the nested representation sent to the frontend.
type FileTreeNode struct {
	Name     string          `json:"name"`
	Path     string          `json:"path"`
	Type     string          `json:"type"`
	Size     int64           `json:"size,omitempty"`
	Language *string         `json:"language,omitempty"`
	Children []*FileTreeNode `json:"children,omitempty"`
}

// FileContent is returned when the frontend requests a specific file's content.
type FileContent struct {
	Path     string  `json:"path"`
	Name     string  `json:"name"`
	Content  string  `json:"content"`
	Language *string `json:"language,omitempty"`
	Size     int64   `json:"size"`
}
