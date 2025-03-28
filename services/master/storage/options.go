package storage

func WithBufferSize(size int) func(worker *Worker) {
	return func(worker *Worker) {
		worker.bufferSize = size
	}
}
