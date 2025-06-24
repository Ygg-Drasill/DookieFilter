package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/data"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/filter"
	"github.com/Ygg-Drasill/DookieFilter/services/master/proxy"
	"github.com/Ygg-Drasill/DookieFilter/services/master/storage"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	"github.com/joho/godotenv"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(logger.New("master", "DEBUG"))
}

const dataWindowSize = 5 * 25 //seconds * frames per seconds

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Failed to load env variables")
		return
	}
	slog.Info("Starting master service")
	socketCtx, err := zmq.NewContext()
	if err != nil {
		slog.Error(err.Error())
	}

	workers := worker.NewPool()

	dataPath := os.Getenv("MATCH_FILE")
	fr, err := frameReader.New(dataPath)
	err = fr.GoToFrame(fr.FrameCount() / 2)
	if err != nil {
		slog.Error("Failed to make frame loader", "error", err.Error())
		return
	}
	dl := data.New(socketCtx, "inproc://data", data.WithFrameLoader(fr))

	workers.Add(dl)

	workers.Add(collector.New(socketCtx, collector.WithEndpoint("inproc://data")))

	workers.Add(storage.New(socketCtx,
		storage.WithBufferSize(dataWindowSize)))

	workers.Add(detector.New(socketCtx))

	workers.Add(filter.New(socketCtx))

	workers.Add(proxy.New(socketCtx))

	workers.Wait()
}
