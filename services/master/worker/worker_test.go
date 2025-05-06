package worker

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBase(t *testing.T) {
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)
	baseWorker := NewBaseWorker(ctx, "test")
	assert.NotNil(t, baseWorker)
	assert.NotNil(t, baseWorker.Logger)
	assert.NotNil(t, baseWorker.SocketContext)
	assert.Equal(t, ctx, baseWorker.SocketContext)
}
