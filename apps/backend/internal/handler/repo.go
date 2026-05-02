package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/internal/worker"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/go-chi/chi/v5"
)

type RepoHandler struct {
	repoService service.RepoServiceInterface
	workerPool  worker.WorkerPoolInterface
}

func NewRepoHandler(repoService service.RepoServiceInterface, workerPool worker.WorkerPoolInterface) *RepoHandler {
	return &RepoHandler{
		repoService: repoService,
		workerPool:  workerPool,
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

	h.workerPool.AddJob(models.Job{
		RepoID:  repo.ID,
		RepoURL: repo.URL,
	})

	resp := models.CreateRepoResponse{
		ID:     repo.ID,
		URL:    repo.URL,
		Status: repo.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *RepoHandler) GetRepo(w http.ResponseWriter, r *http.Request) {
	//get id from url
	//
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	repo, err := h.repoService.GetRepo(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get repo", http.StatusInternalServerError)
		return
	}

	resp := models.GetRepoResponse{
		ID:        repo.ID,
		URL:       repo.URL,
		Status:    repo.Status,
		CreatedAt: repo.CreatedAt,
		UpdatedAt: repo.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RepoHandler) GetRepoDependencies(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	deps, err := h.repoService.GetRepoDependencies(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get dependencies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deps)
}
