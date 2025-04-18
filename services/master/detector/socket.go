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

	err = w.socketListen.Bind(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) close() {
	err := w.socketListen.Close()
	if err != nil {
		w.Logger.Warn("failed to close socket (forward)", "error", err.Error())
	}
}
