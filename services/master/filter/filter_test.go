package filter

import (
	"sync"
	"testing"

	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
)

func Test_FilterWorker(t *testing.T) {
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)
	inputSocket, err := ctx.NewSocket(zmq.PUSH)
	assert.NoError(t, err)
	outputSocket, err := ctx.NewSocket(zmq.PULL)
	assert.NoError(t, err)
	endpoint := "test"
	const framesCount = 10

	worker := New(ctx, WithOutputEndpoint(endpoint))

	wg := &sync.WaitGroup{}
	worker.Run(wg)
	
	inputSocket.Connect(endpoints.InProcessEndpoint(endpoints.FILTER_INPUT))
	outputSocket.Connect(endpoint)
	
	frames := make([]types.Frame,1)
	frames[0] = testutils.RandomFrame(3,3)
}
