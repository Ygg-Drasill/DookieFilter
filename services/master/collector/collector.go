package collector

import (
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/services/master/config"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
	"sync"
)

type CollectorWorker struct {
	logger         *slog.Logger
	socketContext  *zmq.Context
	socketListen   *zmq.Socket
	socketStore    *zmq.Socket
	socketDetector *zmq.Socket
}

func New(options ...func(worker *CollectorWorker)) *CollectorWorker {
	worker := &CollectorWorker{}
	for _, opt := range options {
		opt(worker)
	}
	worker.logger = logger.New(config.ServiceName, config.DebugLevel, slog.Attr{
		Key:   "worker",
		Value: slog.StringValue("collector"),
	})
	return worker
}

func (w *CollectorWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer w.close()
	w.logger.Info("starting collector worker")
	err := w.connect()
	if err != nil {
		w.logger.Error("failed to bind/connect zmq sockets", "error", err.Error())
	}

	for {
		err := w.listen()
		if err != nil {
			w.logger.Error(err.Error())
		}
	}

	w.logger.Warn("collector worker stopped")
}
