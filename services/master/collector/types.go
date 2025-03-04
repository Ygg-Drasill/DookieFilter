package collector

import (
    zmq "github.com/pebbe/zmq4"
    "log/slog"
    "strings"
    "sync"
)

type CollectorWorker struct {
    socketContext    *zmq.Context
    socketListen     *zmq.Socket
    socketStore      *zmq.Socket
    socketDistribute *zmq.Socket
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
    var err error
    w.socketListen, err = w.socketContext.NewSocket(zmq.PULL)
    if err != nil {
        slog.Error(err.Error())
    }

    w.socketStore, err = w.socketContext.NewSocket(zmq.PUSH)
    if err != nil {
        slog.Error(err.Error())
    }

    err = w.socketListen.Bind("inproc://collector")
    if err != nil {
        slog.Error(err.Error())
    }

    for {
        err := w.collectorListen()
        if err != nil {
            slog.Error(err.Error())
        }
    }
}

func (w *CollectorWorker) collectorListen() error {
    topic, err := w.socketListen.Recv(zmq.SNDMORE)
    if err != nil {
        return err
    }

    msg, err := w.socketListen.RecvMessage(0)
    if err != nil {
        return err
    }

    if topic == "frame" {
        //parse raw frame
    }

    if topic == "point" {
        _, err := w.socketStore.SendMessage(0, "point", msg)
        if err != nil {
            return err
        }
    }

    return nil
}
