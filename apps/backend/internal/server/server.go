package server

import (
	"context"
	"log"
	"net/http"

	handlers "github.com/Aneeshie/repo-analyzer/backend/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router      *chi.Mux
	db          *pgxpool.Pool
	port        string
	repoHandler *handlers.RepoHandler
}

func NewServer(db *pgxpool.Pool, repoHandler *handlers.RepoHandler) *Server {
	server := &Server{
		router:      chi.NewRouter(),
		db:          db,
		port:        ":8080",
		repoHandler: repoHandler,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {

	s.router.Get("/health", s.healthCheck)
	s.router.Post("/api/v1/repos", s.repoHandler.CreateRepo)
	s.router.Get("/api/v1/repos/{id}", s.repoHandler.GetRepo)
	s.router.Get("/api/v1/repos/{id}/dependencies", s.repoHandler.GetRepoDependencies)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(context.Background()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("database connection failed"))
		return
	}
	w.Write([]byte("OK"))
}

func (s *Server) Run() error {
	s.setupRoutes()
	log.Printf("Server starting on %s", s.port)
	return http.ListenAndServe(s.port, s.router)
}
