package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
	"github.com/stretchr/testify/assert"
)

// Mock service implementing RepoServiceInterface
type mockRepoService struct {
	createFunc  func(ctx context.Context, url string) (*models.Repo, error)
	getFunc     func(ctx context.Context, id string) (*models.Repo, error)
	updateFunc  func(ctx context.Context, id string, status string) error
	getDepsFunc func(ctx context.Context, id string) ([]models.Dependency, error)
}

func (m *mockRepoService) CreateRepo(ctx context.Context, url string) (*models.Repo, error) {
	if m.createFunc == nil {
		return &models.Repo{ID: "123", URL: url, Status: models.StatusPending}, nil
	}
	return m.createFunc(ctx, url)
}

func (m *mockRepoService) GetRepo(ctx context.Context, id string) (*models.Repo, error) {
	if m.getFunc == nil {
		return &models.Repo{ID: id, URL: "https://github.com/test/repo", Status: models.StatusCompleted}, nil
	}
	return m.getFunc(ctx, id)
}

func (m *mockRepoService) UpdateRepoStatus(ctx context.Context, id string, status string) error {
	if m.updateFunc == nil {
		return nil
	}
	return m.updateFunc(ctx, id, status)
}

func (m *mockRepoService) GetRepoDependencies(ctx context.Context, id string) ([]models.Dependency, error) {
	if m.getDepsFunc == nil {
		return []models.Dependency{}, nil
	}
	return m.getDepsFunc(ctx, id)
}

// Mock worker pool implementing WorkerPoolInterface
type mockWorkerPool struct {
	addJobFunc func(job models.Job)
}

func (m *mockWorkerPool) AddJob(job models.Job) {
	if m.addJobFunc != nil {
		m.addJobFunc(job)
	}
}

func (m *mockWorkerPool) Shutdown() {
	// Do nothing
}

func TestCreateRepo_Success(t *testing.T) {
	mockService := &mockRepoService{}
	mockPool := &mockWorkerPool{}

	handler := NewRepoHandler(mockService, mockPool)

	reqBody := `{"url": "https://github.com/test/repo"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/repos", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.CreateRepo(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.CreateRepoResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/test/repo", resp.URL)
	assert.Equal(t, "pending", resp.Status)
}

func TestCreateRepo_InvalidBody(t *testing.T) {
	mockService := &mockRepoService{}
	mockPool := &mockWorkerPool{}
	handler := NewRepoHandler(mockService, mockPool)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/repos", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRepo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateRepo_MissingURL(t *testing.T) {
	mockService := &mockRepoService{}
	mockPool := &mockWorkerPool{}
	handler := NewRepoHandler(mockService, mockPool)

	reqBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/repos", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRepo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
