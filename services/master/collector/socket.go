package collector

import (
	zmq "github.com/pebbe/zmq4"
)

func (w *CollectorWorker) connect() error {
	var err error
	w.socketListen, err = w.socketContext.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}

	w.socketStore, err = w.socketContext.NewSocket(zmq.PUSH)
	if err != nil {
		return err
	}

	w.socketForward, err = w.socketContext.NewSocket(zmq.PUSH)
	if err != nil {
		return err
	}

	err = w.socketListen.Bind("inproc://collector")
	if err != nil {
		return err
	}

	err = w.socketStore.Bind("inproc://storage")
	if err != nil {
		return err
	}

	err = w.socketForward.Bind("inproc://detector")
	if err != nil {
		return err
	}

	return nil
}

func (w *CollectorWorker) close() {
	var err error
	err = w.socketListen.Close()
	if w.socketStore != nil {
		w.logger.Warn("failed to close socket (listen)", "error", err.Error())
	}
	err = w.socketStore.Close()
	if w.socketStore != nil {
		w.logger.Warn("failed to close socket (storage)", "error", err.Error())
	}
	err = w.socketForward.Close()
	if w.socketStore != nil {
		w.logger.Warn("failed to close socket (forward)", "error", err.Error())
	}

}
