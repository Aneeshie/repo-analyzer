package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type FileTreeHandler struct {
	fileTreeService *service.FileTreeService
	storagePath     string
}

func NewFileTreeHandler(fileTreeService *service.FileTreeService, storagePath string) *FileTreeHandler {
	return &FileTreeHandler{
		fileTreeService: fileTreeService,
		storagePath:     storagePath,
	}
}

// GetFileTree returns the full nested file tree for a repo.
func (h *FileTreeHandler) GetFileTree(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	tree, err := h.fileTreeService.GetTree(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get file tree", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tree)
}

// GetFileContent reads a file's content from disk and returns it.
func (h *FileTreeHandler) GetFileContent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "path query parameter is required", http.StatusBadRequest)
		return
	}

	repoPath := filepath.Join(h.storagePath, id)

	// Make sure the repo directory actually exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.Error(w, "Repository files not found on disk", http.StatusNotFound)
		return
	}

	content, err := h.fileTreeService.ReadFileContent(repoPath, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}
