package filter

func WithOutputEndpoint(endpoint string) func (w *Worker) {
	return func(w *Worker) {
		w.outputEndpoint = endpoint
	}
}
