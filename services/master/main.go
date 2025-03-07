package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/data"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
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

	fr, err := frameReader.New("../../visunator/raw.jsonl")
	err = fr.GoToFrame(fr.FrameCount() / 2)
	if err != nil {
		slog.Error("Failed to make frame loader", "error", err.Error())
		return
	}
	workers.Add(data.New(socketCtx, data.WithFrameLoader(fr)))

	workers.Add(collector.New(socketCtx))

	workers.Add(storage.New(socketCtx,
		storage.WithBufferSize(1024)))

	workers.Add(detector.New(socketCtx))

	workers.Wait()
}
