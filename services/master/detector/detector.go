package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/services/master/worker"
    zmq "github.com/pebbe/zmq4"
    "sync"
)

type DetectorWorker struct {
    worker.BaseWorker

    socketListen *zmq.Socket
}

func New(ctx *zmq.Context, options ...func(worker *DetectorWorker)) *DetectorWorker {
    w := &DetectorWorker{
        BaseWorker: worker.NewBaseWorker(ctx, "detector"),
    }
    for _, opt := range options {
        opt(w)
    }

    return w
}

func (w *DetectorWorker) Run(wg *sync.WaitGroup) {
    defer wg.Done()
    defer w.close()
    w.Logger.Info("Starting detector worker")
    err := w.connect()
    if err != nil {
        w.Logger.Error("Failed to bind/connect zmq sockets", "error", err)
    }

    for {

    }
}
