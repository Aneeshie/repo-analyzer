package worker

import (
	"context"
	"log"
	"sync"

	"github.com/Aneeshie/repo-analyzer/backend/internal/service"
	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type Pool struct {
	jobQueue      chan models.Job
	repoProcessor *RepoProcessor
	workerCount   int

	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewPool(repoService *service.RepoService, githubService *service.GitHubService, storagePath string, workerCount int) *Pool {

	pool := &Pool{
		jobQueue:      make(chan models.Job, 100),
		repoProcessor: NewRepoProcessor(repoService, githubService, storagePath),
		workerCount:   workerCount,

		stopChan: make(chan struct{}),
	}
	pool.start()
	return pool
}

func (p *Pool) start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func(workerId int) {

			defer p.wg.Done()
			log.Printf("Worker %d started", workerId)

			for {
				select {
				case job := <-p.jobQueue:
					log.Printf("Worker %d processing job: %s", workerId, job.RepoID)
					p.repoProcessor.ProcessRepo(context.Background(), job.RepoID, job.RepoURL)

				case <-p.stopChan:
					log.Printf("worker %d stopping", workerId)
					return
				}
			}

		}(i)
	}
}

func (p *Pool) AddJob(job models.Job) {
	p.jobQueue <- job
}

func (p *Pool) Shutdown() {
	log.Println("Shutting down worker pool...")
	close(p.stopChan)
	p.wg.Wait()
	close(p.jobQueue)
	log.Println("Worker pool shutdown complete")
}
