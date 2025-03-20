package data

import (
	"encoding/json"
	"errors"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"io"
	"os"
	"sync"
	"time"
)

type Worker struct {
	worker.BaseWorker
	socketSend  *zmq.Socket
	frameLoader types.FrameLoader[types.Frame]
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
	w := &Worker{
		BaseWorker: worker.NewBaseWorker(ctx, "dataloader"),
	}
	for _, opt := range options {
		opt(w)
	}
	return w
}

func WithFrameLoader(frameLoader types.FrameLoader[types.Frame]) func(worker *Worker) {
	return func(worker *Worker) {
		worker.frameLoader = frameLoader
	}
}

func (w Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	w.connect()
	w.Logger.Info("Starting dataloader worker")
	lastTime := time.Now().UnixMilli()

	for {
		currentTime := time.Now().UnixMilli()
		if (currentTime - lastTime) < 40 {
			continue
		}
		lastTime = currentTime
		frame, err := w.frameLoader.Next()
		if errors.Is(err, io.EOF) {
			w.Logger.Warn("EOF", "error", err.Error())
			continue
		}
		if errors.Is(err, os.ErrClosed) {
			w.Logger.Warn("FrameLoader closed", "error", err.Error())
			break
		}
		if err != nil {
			w.Logger.Error("Failed to read next frame", "error", err.Error())
		}

		message, err := json.Marshal(frame)
		if err != nil {
			w.Logger.Error("Failed to serialize frame")
		}
		messageLength, err := w.socketSend.SendMessage("frame", message)
		if err != nil {
			w.Logger.Error("Failed to send message", "length", messageLength)
		}
	}
}

func (w *Worker) connect() {
	var err error
	w.socketSend, err = w.SocketContext.NewSocket(zmq.PUSH)
	if err != nil {
		w.Logger.Error("Failed to create new socket")
	}

	err = w.socketSend.Connect(endpoints.InProcessEndpoint(endpoints.COLLECTOR))
	if err != nil {
		w.Logger.Error("Failed to connect socket")
	}
}
