package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/config"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
	"sync"
)

type StorageWorker struct {
	logger        *slog.Logger
	socketContext *zmq.Context
	socketProvide *zmq.Socket
	socketConsume *zmq.Socket

	bufferSize int
	players    map[string]pringleBuffer.PringleBuffer[types.PlayerPosition]
}

func New(options ...func(worker *StorageWorker)) *StorageWorker {
	worker := &StorageWorker{
		bufferSize: 10,
		players:    make(map[string]pringleBuffer.PringleBuffer[types.PlayerPosition]),
	}
	for _, opt := range options {
		opt(worker)
	}
	worker.logger = logger.New(config.ServiceName, config.DebugLevel, slog.Attr{
		Key:   "worker",
		Value: slog.StringValue("storage"),
	})
	return worker
}

func (w *StorageWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer w.close()
	err := w.connect()
	w.logger.Info("starting storage worker")
	if err != nil {
		w.logger.Error("failed to bind/connect zmq sockets", "error", err.Error())
	}

	listenerWaitGroup := &sync.WaitGroup{}
	listenerWaitGroup.Add(2)
	go w.listenProvide(listenerWaitGroup)
	go w.listenConsume(listenerWaitGroup)
	listenerWaitGroup.Wait()
	w.logger.Warn("storage worker stopped")
}
