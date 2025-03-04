package collector

import (
	zmq "github.com/pebbe/zmq4"
	"log/slog"
	"sync"
)

type CollectorWorker struct {
	socketContext *zmq.Context
}

func New(options ...func(worker *CollectorWorker)) *CollectorWorker {
	worker := &CollectorWorker{}
	for _, opt := range options {
		opt(worker)
	}
	return worker
}

func (w *CollectorWorker) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	socket, err := w.socketContext.NewSocket(zmq.PUSH)
	if err != nil {
		slog.Error(err.Error())
	}
	err = socket.Bind("inproc://collector")
	slog.Error(err.Error())
}
