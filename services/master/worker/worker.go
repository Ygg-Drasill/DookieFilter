package worker

import (
    "github.com/Ygg-Drasill/DookieFilter/common/logger"
    "github.com/Ygg-Drasill/DookieFilter/services/master/config"
    zmq "github.com/pebbe/zmq4"
    "log/slog"
    "sync"
)

type Worker interface {
    Run(wg *sync.WaitGroup)
}

type BaseWorker struct {
    Logger        *slog.Logger
    SocketContext *zmq.Context
}

func NewBaseWorker(socketContext *zmq.Context, workerName string) BaseWorker {
    return BaseWorker{
        Logger: logger.New(config.ServiceName, config.DebugLevel, slog.Attr{
            Key:   "worker",
            Value: slog.StringValue(workerName),
        }),
        SocketContext: socketContext,
    }
}
