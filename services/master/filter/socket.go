package filter

import (
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	zmq "github.com/pebbe/zmq4"
)

func (w *Worker) connect() error {
	w.socketInput = w.BaseWorker.NewSocket(zmq.PULL)
	w.socketOutput = w.BaseWorker.NewSocket(zmq.PUSH)

	w.BaseWorker.Bind(w.socketInput, endpoints.InProcessEndpoint(endpoints.FILTER_INPUT))
	w.BaseWorker.Bind(w.socketOutput, w.outputEndpoint)

	return nil //TODO: no return value (but for all workers)
}

func (w *Worker) close() {
	w.BaseWorker.CloseAllSockets()
}
