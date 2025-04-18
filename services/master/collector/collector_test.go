package collector

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"sync"
	"testing"
)

func TestCollector(t *testing.T) {

	wg := &sync.WaitGroup{}
	ctx, err := zmq.NewContext()
	ctx.SetIoThreads(8)
	if err != nil {
		panic(err)
	}
	worker := New(ctx)
	wg.Add(1)

	socketCollector, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		panic(err)
	}
	socketStore, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}
	socketDetector, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}

	err = socketStore.Bind(endpoints.InProcessEndpoint(endpoints.STORAGE))
	assert.NoError(t, err)
	err = socketDetector.Bind(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	assert.NoError(t, err)

	go worker.Run(wg)

	err = socketCollector.Connect(endpoints.InProcessEndpoint(endpoints.COLLECTOR))
	assert.NoError(t, err)

	frame := testutils.RandomFrame(11, 11)
	_, err = json.Marshal(frame)
	if err != nil {
		panic(err)
	}
	go func() {
		n, err := socketCollector.SendMessage("frame")
		assert.Greater(t, n, 0)
		if err != nil {
			panic(err)
		}
	}()

	//Wrapping this part of the test in a separate goroutine
	//allows the zmq sockets to run in parallel.
	//Otherwise, socketStore.RecvMessage halts execution for the sockets
	//inside the worker, which would not be able to receive messages.
	//(Running it in a t.Run() with t.Parallel() did not work...)
	testWg := &sync.WaitGroup{}
	var storeMessage []string
	var detectorMessage []string
	testWg.Add(1)
	go func() {
		storeMessage, err = socketStore.RecvMessage(0)
		assert.NoError(t, err)
		slog.Info("hello")

		detectorMessage, err = socketDetector.RecvMessage(0)
		assert.NoError(t, err)
		testWg.Done()
	}()
	testWg.Wait()
	assert.NotEmpty(t, storeMessage)
	assert.NotEmpty(t, detectorMessage)
}
