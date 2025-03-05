package worker

import (
	"sync"
)

type workerPool struct {
	workers   []Worker
	waitGroup sync.WaitGroup
}

func NewPool() *workerPool {
	return &workerPool{
		workers:   make([]Worker, 0),
		waitGroup: sync.WaitGroup{},
	}
}

func (pool *workerPool) Add(worker Worker) {
	pool.waitGroup.Add(1)
	pool.workers = append(pool.workers, worker)
	go worker.Run(&pool.waitGroup)
}

func (pool *workerPool) Wait() {
	pool.waitGroup.Wait()
}
