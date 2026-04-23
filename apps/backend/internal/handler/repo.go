package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type RepoHandler struct {
	repoService *service.RepoService
}

func NewRepoHandler(repoService *service.RepoService) *RepoHandler {
	return &RepoHandler{
		repoService: repoService,
	}
}

func (h *RepoHandler) CreateRepo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRepoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	repo, err := h.repoService.CreateRepo(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "Failed to create repo", http.StatusInternalServerError)
		return
	}

	resp := models.CreateRepoResponse{
		ID:     repo.ID,
		URL:    repo.URL,
		Status: repo.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
