package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/go-chi/chi/v5"
)

type DependencyHandler struct{
	depService *service.DependencyService
}

func NewDependencyHandler(depService *service.DependencyService) *DependencyHandler {
	return &DependencyHandler{
		depService: depService,
	}
}

func (h *DependencyHandler) GetDependencies(w http.ResponseWriter, r *http.Request){
	repoID := chi.URLParam(r, "id")
	if repoID == ""{
		http.Error(w, "Repo ID is required", http.StatusBadRequest)
		return
	}

	deps, err := h.depService.GetDependencyByID(r.Context(), repoID)
	if err != nil {
		http.Error(w, "Failed to fetch dependencies", http.StatusInternalServerError)
		return
	}

	//group by ecosystem
	grouped := make(map[string][]models.Dependency)
	for _, dep := range deps {
		grouped[string(dep.Ecosystem)] = append(grouped[string(dep.Ecosystem)],dep)
	}

	response := struct {
		RepoID       string                             `json:"repo_id"`
		Total        int                                `json:"total"`
		Dependencies map[string][]models.Dependency     `json:"dependencies"`
	}{
		RepoID:       repoID,
		Total:        len(deps),
		Dependencies: grouped,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}


