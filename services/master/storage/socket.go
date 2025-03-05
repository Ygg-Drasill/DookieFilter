package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	zmq "github.com/pebbe/zmq4"
)

func (w *StorageWorker) connect() error {
	var err error
	w.socketConsume, err = w.SocketContext.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}

	w.socketProvide, err = w.SocketContext.NewSocket(zmq.REP)
	if err != nil {
		return err
	}

	err = w.socketConsume.Bind(endpoints.InProcessEndpoint(endpoints.STORAGE))
	if err != nil {
		return err
	}

	err = w.socketProvide.Bind(endpoints.InProcessEndpoint(endpoints.STORAGE_PROVIDE))
	if err != nil {
		return err
	}

	return nil
}

func (w *StorageWorker) close() {
	var err error
	err = w.socketConsume.Close()
	if err != nil {
		w.Logger.Error("failed to close socket", "error", err.Error())
	}
	err = w.socketProvide.Close()
	if err != nil {
		w.Logger.Error("failed to close socket", "error", err.Error())
	}
}
