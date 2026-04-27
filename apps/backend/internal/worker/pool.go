package worker

import (
	"context"
	"log"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type Pool struct {
	jobQueue      chan models.Job
	repoProcessor *RepoProcessor
	workerCount   int
}

func NewPool(repoService *service.RepoService, githubService *service.GitHubService, storagePath string, workerCount int) *Pool {

	pool := &Pool{
		jobQueue:      make(chan models.Job, 100),
		repoProcessor: NewRepoProcessor(repoService, githubService, storagePath),
		workerCount:   workerCount,
	}
	pool.start()
	return pool
}

func (p *Pool) start() {
	for i := 0; i < p.workerCount; i++ {
		go func(workerId int) {
			log.Printf("Worker %d started", workerId)
			for job := range p.jobQueue {
				p.repoProcessor.ProcessRepo(context.Background(), job.RepoID, job.RepoURL)
			}
		}(i)
	}
}

func (p *Pool) AddJob(job models.Job) {
	p.jobQueue <- job
}

func (p *Pool) Shutdown() {
	close(p.jobQueue)
}
