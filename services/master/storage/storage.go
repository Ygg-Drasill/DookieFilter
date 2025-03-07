package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

type Worker struct {
	worker.BaseWorker
	socketProvide *zmq.Socket
	socketConsume *zmq.Socket

	bufferSize int
	players    map[string]pringleBuffer.PringleBuffer[types.PlayerPosition]
	mutex      sync.Mutex
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
	w := &Worker{
		BaseWorker: worker.NewBaseWorker(ctx, "storage"),
		bufferSize: 10,
		players:    make(map[string]pringleBuffer.PringleBuffer[types.PlayerPosition]),
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
	w.Logger.Info("Starting storage worker")
	if err != nil {
		w.Logger.Error("Failed to bind/connect zmq sockets", "error", err.Error())
	}

	listenerWaitGroup := &sync.WaitGroup{}
	listenerWaitGroup.Add(2)
	go w.listenProvide(listenerWaitGroup)
	go w.listenConsume(listenerWaitGroup)
	listenerWaitGroup.Wait()
	w.Logger.Warn("Storage worker stopped")
}
