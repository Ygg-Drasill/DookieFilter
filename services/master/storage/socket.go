package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	zmq "github.com/pebbe/zmq4"
)

func (w *Worker) connect() error {
	w.socketConsume = w.BaseWorker.NewSocket(zmq.PULL)
	w.socketProvide = w.BaseWorker.NewSocket(zmq.REP)
	w.socketAPI = w.BaseWorker.NewSocket(zmq.REP)

	w.BaseWorker.Bind(w.socketConsume, endpoints.InProcessEndpoint(endpoints.STORAGE))
	w.BaseWorker.Bind(w.socketProvide, endpoints.InProcessEndpoint(endpoints.STORAGE_PROVIDE))
	w.BaseWorker.Bind(w.socketAPI, endpoints.TcpEndpoint(endpoints.STORAGE_API))
	return nil //TODO: no return value (but for all workers)
}

func (w *Worker) close() {
	w.BaseWorker.CloseAllSockets()
}
