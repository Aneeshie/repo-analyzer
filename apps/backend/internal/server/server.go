package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router *chi.Mux
	db     *pgxpool.Pool
	port   string
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		router: chi.NewRouter(),
		db:     db,
		port:   ":8080",
	}
}

func (s *Server) setupRoutes() {
	s.router.Get("/health", s.healthCheck)
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
