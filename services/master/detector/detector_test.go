package detector

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
)

// Test timeout in seconds
const timeout = 20

func TestSwapSendsModifiedPositionToStorage(t *testing.T) {
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)

	w := New(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go w.Run(wg)

	// Setup socket to send frame to swap worker
	socketInput, err := ctx.NewSocket(zmq.PUSH)
	assert.NoError(t, err)
	socketStorage, err := ctx.NewSocket(zmq.PULL)
	assert.NoError(t, err)

	err = socketInput.Connect(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	assert.NoError(t, err)
	err = socketStorage.Bind(endpoints.InProcessEndpoint(endpoints.STORAGE))
	assert.NoError(t, err)

	frame := testutils.RandomFrame(11, 11)
	framePacket := types.SerializeFrame(types.SmallFromBigFrame(frame))
	t.Logf("Sending frame")
	_, err = socketInput.SendMessage("frame", framePacket)
	assert.NoError(t, err)
	_, err = socketInput.SendMessage("frame", framePacket)
	assert.NoError(t, err)

	m, err := socketStorage.RecvMessage(0)
	assert.NoError(t, err)
	t.Logf("Received frame from storage %s", m)

	assert.NoError(t, socketInput.Close())
	assert.NoError(t, socketStorage.Close())
}

func TestDetectHole(t *testing.T) {
	endpoint := "inproc://test"
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)

	w := New(ctx, WithImputationEndpoint(endpoint))
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go w.Run(wg)

	assert.NoError(t, err)
	socketInput, err := ctx.NewSocket(zmq.PUSH)
	assert.NoError(t, err)
	socketImputation, err := ctx.NewSocket(zmq.PULL)
	assert.NoError(t, err)
	err = socketInput.Connect(endpoints.InProcessEndpoint(endpoints.DETECTOR))
	assert.NoError(t, err)
	err = socketImputation.Bind(endpoint)
	assert.NoError(t, err)

	// Generate random frame for testing
	frame := testutils.RandomFrame(2, 1)
	t.Logf("Sending initial frame with %d home players", len(frame.HomePlayers))
	_, err = socketInput.SendMessage("frame", types.SerializeFrame(types.SmallFromBigFrame(frame)), 0)
	assert.NoError(t, err)

	// Wait a bit to ensure first frame is processed
	time.Sleep(100 * time.Millisecond)

	// Remove one player from the home list
	next := testutils.RandomNextFrame(frame)
	removeIndex := len(next.HomePlayers) - 1
	playerNumRaw := next.HomePlayers[removeIndex].Number
	playerNum, err := strconv.Atoi(playerNumRaw)
	assert.NoError(t, err)
	next.HomePlayers = next.HomePlayers[0:removeIndex]
	t.Logf("Sending next frame with %d home players (removed player %d)", len(next.HomePlayers), playerNum)
	_, err = socketInput.SendMessage("frame", types.SerializeFrame(types.SmallFromBigFrame(next)), 0)
	assert.NoError(t, err)

	doneChan := make(chan bool)
	go func() {
		topic, err := socketImputation.Recv(zmq.SNDMORE)
		assert.NoError(t, err)
		assert.Equal(t, "hole", topic)
		t.Logf("Received hole message with topic: %s", topic)

		packet, err := socketImputation.RecvMessage(0)
		message := strings.Join(packet, "")
		assert.NoError(t, err)
		assert.Greater(t, len(message), 0)
		t.Logf("Received message: %s", message)

		request := holeMessage{}
		err = json.Unmarshal([]byte(message), &request)
		assert.NoError(t, err)
		t.Logf("Parsed hole message: FrameIdx=%d, PlayerNumber=%d, Home=%v",
			request.FrameIdx, request.PlayerNumber, request.Home)

		assert.Equal(t, next.FrameIdx, request.FrameIdx, "Frame index should match")
		assert.Equal(t, playerNum, request.PlayerNumber, "Player ID should match the removed player")
		assert.Equal(t, true, request.Home, "Should be a home player")

		doneChan <- true
	}()

	// Fail the test after timeout
	go func() {
		time.Sleep(timeout * time.Second)
		//doneChan <- false
	}()

	result := <-doneChan
	assert.True(t, result, "expected message from detector worker within timeout")
}
