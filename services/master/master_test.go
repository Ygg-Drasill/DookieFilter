package main

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/storage"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

const testFrameCount = 10

func TestWorkerIntegration(t *testing.T) {
	ctx, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}

	pool := worker.NewPool()
	pool.Add(storage.New(ctx))
	pool.Add(collector.New(ctx))
	pool.Add(detector.New(ctx))

	collectorSock, err := ctx.NewSocket(zmq.PUSH)
	if err != nil {
		panic(err)
	}
	err = collectorSock.Connect(endpoints.InProcessEndpoint(endpoints.COLLECTOR))
	if err != nil {
		panic(err)
	}

	storageSock, err := ctx.NewSocket(zmq.REQ)
	if err != nil {
		panic(err)
	}

	err = storageSock.Connect(endpoints.InProcessEndpoint(endpoints.STORAGE_PROVIDE))
	if err != nil {
		panic(err)
	}

	initialFrame := testutils.RandomFrame(3, 3)
	testFrames := []types.SmallFrame{types.SmallFromBigFrame(initialFrame)}
	for i, frame := 0, initialFrame; i < testFrameCount; i, frame = i+1, testutils.RandomNextFrame(frame) {
		testFrames = append(testFrames, types.SmallFromBigFrame(frame))
		_, err = collectorSock.Send("frame", zmq.SNDMORE)
		if err != nil {
			panic(err)
		}

		frameMarshal, err := json.Marshal(frame)
		if err != nil {
			panic(err)
		}
		_, err = collectorSock.Send(string(frameMarshal), 0)
		assert.NoError(t, err)
	}

	playerIds := make([]string, len(initialFrame.AwayPlayers)+len(initialFrame.HomePlayers))
	for _, p := range append(initialFrame.AwayPlayers, initialFrame.HomePlayers...) {
		playerIds = append(playerIds, p.PlayerId)
	}

	for _, f := range testFrames {
		for _, p := range f.Players {
			var message string
			for message == "" {
				_, err = storageSock.SendMessage("playerFrame", f.FrameIdx, ":", p.PlayerNum)
				if err != nil {
					panic(err)
				}

				parts, err := storageSock.RecvMessage(0)
				if err != nil {
					panic(err)
				}
				message = strings.Join(parts, "")
			}
			rawCoords := strings.Split(message, ";")
			assert.Equal(t, 2, len(rawCoords), "Expected 2 elements in coordinates")
			x, err := strconv.ParseFloat(rawCoords[0], 64)
			assert.NoError(t, err)
			y, err := strconv.ParseFloat(rawCoords[1], 64)
			assert.NoError(t, err)
			assert.Equal(t, x, p.X, "Expected player x to match")
			assert.Equal(t, y, p.Y, "Expected player y to match")
		}
	}
}
