package collector

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCollector(t *testing.T) {
	wg := &sync.WaitGroup{}
	ctx, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}
	worker := New(ctx)
	go worker.Run(wg)

	socketCollector, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		panic(err)
	}
	err = socketCollector.Connect(endpoints.InProcessEndpoint(endpoints.COLLECTOR))
	assert.NoError(t, err)
	socketStore, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}
	err = socketStore.Bind(endpoints.InProcessEndpoint(endpoints.STORAGE))
	assert.NoError(t, err)
	socketDetector, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}
	err = socketDetector.Bind(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	assert.NoError(t, err)

	frame := testutils.RandomFrame(11, 11)
	data, err := json.Marshal(frame)
	if err != nil {
		panic(err)
	}
	_, err = socketCollector.SendMessage("frame", string(data))
	if err != nil {
		panic(err)
	}

	storeMessage, err := socketStore.RecvMessage(0)
	assert.NoError(t, err)
	assert.NotEmpty(t, storeMessage)
	detectorMessage, err := socketDetector.RecvMessage(0)
	assert.NoError(t, err)
	assert.NotEmpty(t, detectorMessage)
}
