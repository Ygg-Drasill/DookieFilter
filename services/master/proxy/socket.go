package proxy

import (
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	zmq "github.com/pebbe/zmq4"
)

func (w *Worker) connect() error {
	w.socketListen = w.BaseWorker.NewSocket(zmq.PULL)
	w.socketForward = w.BaseWorker.NewSocket(zmq.PUSH)

	w.BaseWorker.Bind(w.socketListen, w.endpoint)
	w.BaseWorker.Connect(w.socketForward, endpoints.InProcessEndpoint(endpoints.STORAGE))
	return nil
}

func (w *Worker) close() {
	w.BaseWorker.CloseAllSockets()
}
