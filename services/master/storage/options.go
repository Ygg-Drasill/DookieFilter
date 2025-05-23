package storage

// WithBufferSize sets the size of the storage buffer, default is 10
func WithBufferSize(size int) func(worker *Worker) {
	return func(worker *Worker) {
		worker.bufferSize = size
	}
}

// WithAPIEndpoint sets the API endpoint to an alternative endpoint
func WithAPIEndpoint(endpoint string) func(worker *Worker) {
	return func(worker *Worker) {
		worker.socketAPIAddress = endpoint
	}
}
