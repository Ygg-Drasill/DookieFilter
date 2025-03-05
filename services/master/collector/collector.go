package collector

import (
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

type CollectorWorker struct {
	worker.BaseWorker
	socketListen   *zmq.Socket
	socketStore    *zmq.Socket
	socketDetector *zmq.Socket
}

func New(ctx *zmq.Context, options ...func(worker *CollectorWorker)) *CollectorWorker {
	w := &CollectorWorker{
		BaseWorker: worker.NewBaseWorker(ctx, "collector"),
	}
	for _, opt := range options {
		opt(w)
	}
	return w
}

func (w *CollectorWorker) Run(wg *sync.WaitGroup) {
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
