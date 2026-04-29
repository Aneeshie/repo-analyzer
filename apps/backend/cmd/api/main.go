package main

import (
	"context"
	"log"
	"os"

	handlers "github.com/Aneeshie/repo-analyzer/backend/internal/handler"
	"github.com/Aneeshie/repo-analyzer/backend/internal/repository"
	"github.com/Aneeshie/repo-analyzer/backend/internal/server"
	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No. env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL with connection pool")

	repoRepo := repository.NewRepoRepository(pool)
	repoService := service.NewRepoService(repoRepo)
	githubService := service.NewGitHubService()

	depRepo := repository.NewDependencyRepository(pool)
	depService := service.NewDependencyService(depRepo)
	depHandler := handlers.NewDependencyHandler(depService)

	//get the storage path

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "../storage/repos"
	}

	//create workerPool
	workerPool := worker.NewPool(repoService, githubService, storagePath, pool, 4)
	defer workerPool.Shutdown()

	// pass worker pool to the handler
	repoHandler := handlers.NewRepoHandler(repoService, workerPool)

	srv := server.NewServer(pool, repoHandler, depHandler)
	if err := srv.Run(); err != nil {
		log.Fatal("Server error:", err)
	}

}
