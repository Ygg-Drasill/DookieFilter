package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/storage"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
)

func init() {
	slog.SetDefault(logger.New("master", "DEBUG"))
}

const dataWindowSize = 30 * 25 //seconds * frames per seconds

func main() {
	slog.Info("Starting master service")
	socketCtx, err := zmq.NewContext()
	if err != nil {
		slog.Error(err.Error())
	}

	workers := worker.NewPool()

	workers.Add(collector.New(socketCtx))

	workers.Add(storage.New(socketCtx,
		storage.WithBufferSize(1024)))

	workers.Wait()
}
