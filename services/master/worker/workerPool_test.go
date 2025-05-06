package worker

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type TestWorker struct {
	BaseWorker
	hasRun *bool
}

func NewTestWorker(ctx *zmq.Context) TestWorker {
	return TestWorker{
		BaseWorker: NewBaseWorker(ctx, "test"),
		hasRun:     new(bool),
	}
}

func (worker TestWorker) Run(wg *sync.WaitGroup) {
	*worker.hasRun = true
	wg.Done()
}

func TestWorkerPool(t *testing.T) {
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)
	pool := NewPool()
	worker := NewTestWorker(ctx)
	pool.Add(worker)
	pool.Wait()

	assert.True(t, *worker.hasRun)
}
