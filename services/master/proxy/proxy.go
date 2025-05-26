package collector

import (
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

type Worker struct {
	worker.BaseWorker
	socketListen   *zmq.Socket
	socketStore    *zmq.Socket
	socketDetector *zmq.Socket
	endpoint       string
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
	w := &Worker{
		BaseWorker: worker.NewBaseWorker(ctx, "collector"),
		endpoint:   endpoints.TcpEndpoint(endpoints.COLLECTOR),
	}
	for _, opt := range options {
		opt(w)
	}
	return w
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer w.close()
	w.Logger.Info("Starting collector worker")
	err := w.connect()
	if err != nil {
		w.Logger.Error("Failed to bind/connect zmq sockets", "error", err.Error())
	}

	for {
		err := w.listen()
		if err != nil {
			w.Logger.Error(err.Error())
		}
	}

	w.Logger.Warn("Collector worker stopped")
}
