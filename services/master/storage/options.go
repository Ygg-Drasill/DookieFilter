package storage

func WithBufferSize(size int) func(worker *StorageWorker) {
	return func(worker *StorageWorker) {
		worker.bufferSize = size
	}
}
