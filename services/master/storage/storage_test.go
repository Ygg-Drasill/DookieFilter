package storage

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStorageWorker(t *testing.T) {
	const frameCount = 50
	const teamSize = 11
	const endpoint = "inproc://test"
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)
	wg := &sync.WaitGroup{}

	worker := New(ctx, WithBufferSize(frameCount), WithAPIEndpoint(endpoint))
	wg.Add(1)
	go worker.Run(wg)

	time.Sleep(time.Second)

	inputSocket, err := ctx.NewSocket(zmq.PUSH)
	assert.NoError(t, err)
	outputSocket, err := ctx.NewSocket(zmq.REQ)
	assert.NoError(t, err)

	assert.NoError(t, inputSocket.Connect(endpoints.InProcessEndpoint(endpoints.STORAGE)))
	assert.NoError(t, outputSocket.Connect(endpoint))

	frames := testutils.RandomFrameRange(teamSize, frameCount)
	for _, frame := range frames {
		n, err := inputSocket.SendMessage("frame", types.SerializeFrame(types.SmallFromBigFrame(frame)))
		assert.Greater(t, n, 0)
		assert.NoError(t, err)
	}

	time.Sleep(time.Millisecond * 1000)

	assert.Equal(t, teamSize*2, len(worker.players))
	expectedSize := int(math.Max(float64(worker.bufferSize), float64(frameCount)))
	playerKeys := make([]types.PlayerKey, len(worker.players))
	for pNum, buffer := range worker.players {
		assert.Equal(t, expectedSize, buffer.Count(), "player buffer %d should have size %d", pNum, expectedSize)
		playerKeys = append(playerKeys, pNum)
	}

	for _, key := range playerKeys {
		number := key.PlayerNumber
		assert.NoError(t, err, "Player number should be a number")
		request := FrameRangeRequest{
			StartIndex:   frames[0].FrameIdx,
			EndIndex:     frames[len(frames)-1].FrameIdx,
			NearestCount: 3,
			TargetPlayer: number,
		}
		requestData, err := json.Marshal(request)
		assert.NoError(t, err)
		n, err := outputSocket.SendMessage(REQ_FRAME_RANGE, requestData)
		assert.Greater(t, n, 0)
		assert.NoError(t, err)

		resFrames := make([]types.SmallFrame, 0)
		responseData, err := outputSocket.RecvMessage(0)
		assert.NoError(t, err)
		response := strings.Join(responseData, "")
		assert.NoError(t, json.Unmarshal([]byte(response), &resFrames))
		assert.Equal(t, frameCount, len(resFrames))
	}

	wg.Done()
}
