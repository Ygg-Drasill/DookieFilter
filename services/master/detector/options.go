package detector

func WithImputationEndpoint(endpoint string) func(w *Worker) {
	return func(w *Worker) {
		w.imputationEndpoint = endpoint
	}
}
