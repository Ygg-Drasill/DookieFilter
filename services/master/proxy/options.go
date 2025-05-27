package proxy

func WithEndpoint(endpoint string) func(worker *Worker) {
	return func(worker *Worker) {
		worker.endpoint = endpoint
	}
}
