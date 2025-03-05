package detector

import (
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

type DetectorWorker struct {
	worker.BaseWorker
}

func New(ctx *zmq.Context, options ...func(worker *DetectorWorker)) *DetectorWorker {
	w := &DetectorWorker{
		worker.NewBaseWorker(ctx, "detector"),
	}
	for _, opt := range options {
		opt(w)
	}

	return w
}

func (d DetectorWorker) GetBaseWorker() *worker.BaseWorker {
	return &d.BaseWorker
}

func (d DetectorWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	//TODO implement me
	panic("implement me")
}
