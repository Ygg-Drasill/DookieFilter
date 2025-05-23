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
	Sockets       []*zmq.Socket
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

func (w *BaseWorker) NewSocket(t zmq.Type) *zmq.Socket {
	socket, err := w.SocketContext.NewSocket(t)
	if err != nil {
		w.Logger.Error("Failed to create socket",
			"type", t.String(), "error", err.Error())
	}
	return socket
}

func (w *BaseWorker) Connect(socket *zmq.Socket, endpoint string) {
	err := socket.Connect(endpoint)
	if err != nil {
		w.Logger.Error("Failed to connect socket",
			"endpoint", endpoint, "err", err)
	}
}

func (w *BaseWorker) Bind(socket *zmq.Socket, endpoint string) {
	err := socket.Bind(endpoint)
	if err != nil {
		w.Logger.Error("Failed to bind socket",
			"endpoint", endpoint, "err", err)
	}
}

func (w *BaseWorker) CloseAllSockets() {
	for _, socket := range w.Sockets {
		err := socket.Close()
		if err != nil {
			w.Logger.Error("Failed to close socket",
				"socket", socket, "err", err)
		}
	}
}
