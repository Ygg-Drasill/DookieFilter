package storage

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestStorageWorker(t *testing.T) {
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)
	wg := &sync.WaitGroup{}

	worker := New(ctx)
	wg.Add(1)
	go worker.Run(wg)

	wg.Done()
}
