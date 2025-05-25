package collector

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
)

func TestCollector(t *testing.T) {
	endpoint := "inproc://test"
	workerWg := &sync.WaitGroup{}
	ctx, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}
	worker := New(ctx, WithEndpoint(endpoint))

	socketCollector, err := ctx.NewSocket(zmq.PUSH)
	assert.NoError(t, err)
	socketStore, err := ctx.NewSocket(zmq.PULL)
	assert.NoError(t, err)
	socketDetector, err := ctx.NewSocket(zmq.PULL)
	assert.NoError(t, err)

	err = socketStore.Bind(endpoints.InProcessEndpoint(endpoints.STORAGE))
	assert.NoError(t, err)
	err = socketDetector.Bind(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	assert.NoError(t, err)
	err = socketCollector.Connect(endpoint)
	assert.NoError(t, err)

	workerWg.Add(1)
	go worker.Run(workerWg)

	playerCount := 11
	frame := testutils.RandomFrame(playerCount, playerCount)
	framePacket, err := json.Marshal(frame)
	if err != nil {
		panic(err)
	}
	n, err := socketCollector.SendMessage("frame", framePacket)
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	var storeMessage []string
	var detectorMessage []string

	storeMessage, err = socketStore.RecvMessage(0)
	assert.NoError(t, err)
	detectorMessage, err = socketDetector.RecvMessage(0)
	assert.NoError(t, err)

	assert.NotEmpty(t, storeMessage)
	assert.NotEmpty(t, detectorMessage)

	storeFrame := types.DeserializeFrame(strings.Join(storeMessage, ""))
	detectorFrame := types.DeserializeFrame(strings.Join(detectorMessage, ""))

	assert.Len(t, storeFrame.Players, playerCount*2)
	assert.Len(t, detectorFrame.Players, playerCount*2)

	assert.NoError(t, socketCollector.Close())
	assert.NoError(t, socketStore.Close())
	assert.NoError(t, socketDetector.Close())
}
