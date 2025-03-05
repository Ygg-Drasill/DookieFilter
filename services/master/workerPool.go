package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"sync"
)

type workerPool struct {
	workers   []types.Worker
	waitGroup sync.WaitGroup
}

func newWorkerPool() *workerPool {
	return &workerPool{
		workers:   make([]types.Worker, 0),
		waitGroup: sync.WaitGroup{},
	}
}

func (pool *workerPool) Add(worker types.Worker) {
	pool.waitGroup.Add(1)
	pool.workers = append(pool.workers, worker)
	go worker.Run(&pool.waitGroup)
}

func (pool *workerPool) Wait() {
	pool.waitGroup.Wait()
}
