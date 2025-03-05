package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/storage"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
)

func init() {
	slog.SetDefault(logger.New("master", "DEBUG"))
}

const dataWindowSize = 30 * 25 //seconds * frames per seconds

func main() {
	slog.Info("starting master service")
	socketCtx, err := zmq.NewContext()
	if err != nil {
		slog.Error(err.Error())
	}

	workers := newWorkerPool()

	workers.Add(collector.New(
		collector.WithSocketContext(socketCtx)))

	workers.Add(storage.New(
		storage.WithSocketContext(socketCtx),
		storage.WithBufferSize(1024)))

	workers.Wait()
}
