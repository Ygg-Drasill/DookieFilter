package worker

import (
    "sync"
)

type Pool struct {
    workers   []Worker
    waitGroup sync.WaitGroup
}

func NewPool() *Pool {
    return &Pool{
        workers:   make([]Worker, 0),
        waitGroup: sync.WaitGroup{},
    }
}

func (pool *Pool) Add(worker Worker) {
    pool.waitGroup.Add(1)
    pool.workers = append(pool.workers, worker)
    go worker.Run(&pool.waitGroup)
}

func (pool *Pool) Wait() {
    pool.waitGroup.Wait()
}
