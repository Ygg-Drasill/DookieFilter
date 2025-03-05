package storage

import zmq "github.com/pebbe/zmq4"

func WithSocketContext(ctx *zmq.Context) func(worker *StorageWorker) {
	return func(worker *StorageWorker) {
		worker.socketContext = ctx
	}
}

func WithBufferSize(size int) func(worker *StorageWorker) {
	return func(worker *StorageWorker) {
		worker.bufferSize = size
	}
}
