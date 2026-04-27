package worker

import "github.com/Aneeshie/repo-analyzer/backend/pkg/models"

type WorkerPoolInterface interface {
	AddJob(job models.Job)
	Shutdown()
}
