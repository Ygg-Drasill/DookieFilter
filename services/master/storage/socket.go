package storage

import (
	zmq "github.com/pebbe/zmq4"
)

func (w *StorageWorker) connect() error {
	var err error
	w.socketConsume, err = w.socketContext.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}

	w.socketProvide, err = w.socketContext.NewSocket(zmq.REP)
	if err != nil {
		return err
	}

	return nil
}

func (w *StorageWorker) close() {
	var err error
	err = w.socketConsume.Close()
	if err != nil {
		w.logger.Error("failed to close socket", "error", err.Error())
	}
	err = w.socketProvide.Close()
	if err != nil {
		w.logger.Error("failed to close socket", "error", err.Error())
	}
}
