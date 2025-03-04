package collector

import zmq "github.com/pebbe/zmq4"

func WithSocketContext(ctx *zmq.Context) func(worker *CollectorWorker) {
	return func(worker *CollectorWorker) {
		worker.socketContext = ctx
	}
}
