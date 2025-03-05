package storage

import "sync"

func (w *StorageWorker) listenConsume(wg *sync.WaitGroup) {
    defer wg.Done()
}

func (w *StorageWorker) listenProvide(wg *sync.WaitGroup) {
    defer wg.Done()
}
