package detector

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
	w.socketImputation, err = w.SocketContext.NewSocket(zmq.PUSH)
	if err != nil {
		return err
	}

	w.socketStorage, err = w.SocketContext.NewSocket(zmq.PUSH)
	if err != nil {
		return err
	}

	err = w.socketListen.Bind(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	if err != nil {
		return err
	}

	err = w.socketStorage.Connect(endpoints.InProcessEndpoint(endpoints.STORAGE))
	if err != nil {
		return err
	}

	err = w.socketImputation.Connect("tcp://127.0.0.1:5555")
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) close() {
	err := w.socketListen.Close()
	if err != nil {
		w.Logger.Warn("failed to close socket (listen)", "error", err.Error())
	}

	err = w.socketImputation.Close()
	if err != nil {
		w.Logger.Warn("failed to close socket (imputation)", "error", err.Error())
	}
}
