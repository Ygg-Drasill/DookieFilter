package collector

import (
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	zmq "github.com/pebbe/zmq4"
)

func (w *Worker) connect() error {
	var err error
	w.socketListen, err = w.SocketContext.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}

	w.socketStore, err = w.SocketContext.NewSocket(zmq.PUSH)
	if err != nil {
		return err
	}

	w.socketDetector, err = w.SocketContext.NewSocket(zmq.PUSH)
	if err != nil {
		return err
	}

	err = w.socketListen.Bind(endpoints.InProcessEndpoint(endpoints.COLLECTOR))
	if err != nil {
		return err
	}

	err = w.socketStore.Connect(endpoints.InProcessEndpoint(endpoints.STORAGE))
	if err != nil {
		return err
	}

	err = w.socketDetector.Connect(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) close() {
	var err error
	err = w.socketListen.Close()
	if err != nil {
		w.Logger.Warn("failed to close socket (listen)", "error", err.Error())
	}
	err = w.socketStore.Close()
	if err != nil {
		w.Logger.Warn("failed to close socket (storage)", "error", err.Error())
	}
	err = w.socketDetector.Close()
	if err != nil {
		w.Logger.Warn("failed to close socket (forward)", "error", err.Error())
	}
}
