package main

import (
	"context"
	"log"
	"os"

	"github.com/Aneeshie/repo-analyzer/backend/internal/server"
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

	srv := server.NewServer(pool)
	if err := srv.Run(); err != nil {
		log.Fatal("Server error:", err)
	}

}
