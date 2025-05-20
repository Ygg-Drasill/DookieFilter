package filter

import (
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

type Worker struct {
	worker.BaseWorker
	socketInput  *zmq.Socket
	socketOutput *zmq.Socket

	outputEndpoint string
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
	w := &Worker{
		BaseWorker:     worker.NewBaseWorker(ctx, "filter"),
		outputEndpoint: endpoints.TcpEndpoint(endpoints.FILTER_OUTPUT),
	}
	for _, opt := range options {
		opt(w)
	}
	return w
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer w.close()
	err := w.connect()
	w.Logger.Info("Starting filter worker")
	if err != nil {
		w.Logger.Error("Failed to bind/connect zmq sockets", "error", err.Error())
	}

	w.Logger.Warn("Filter worker stopped")
}
