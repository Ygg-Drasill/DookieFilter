package main

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/filter"
	"github.com/Ygg-Drasill/DookieFilter/services/master/storage"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

const testFrameCount = 30

func Test_WorkerIntegration(t *testing.T) {
	ctx, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}

	filterEndpoint := "inproc://filter_test"
	imputationEndpoint := "inproc://imputation_test"
	storageEndpoint := "inproc://storage_test"

	pool := worker.NewPool()
	pool.Add(storage.New(ctx,
		storage.WithBufferSize(testFrameCount), //Make sure that all frames fit in the storage buffer
		storage.WithAPIEndpoint(storageEndpoint)))
	pool.Add(collector.New(ctx))
	pool.Add(detector.New(ctx, detector.WithImputationEndpoint(imputationEndpoint)))
	pool.Add(filter.New(ctx, filter.WithOutputEndpoint(filterEndpoint)))

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

	filterSocket, err := ctx.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}

	err = filterSocket.Connect(filterEndpoint)
	if err != nil {
		panic(err)
	}

	initialFrame := testutils.RandomFrame(3, 3)
	initialFrame.FrameIdx = 1
	testFrames := make([]types.SmallFrame, 0)
	var frame types.Frame
	var i int
	for i, frame = 0, initialFrame; i < testFrameCount; i, frame = i+1, testutils.RandomNextFrame(frame) {
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
	time.Sleep(1 * time.Second)

	for _, f := range testFrames {
		for _, p := range f.Players {
			var message string
			for message == "" {
				_, err = storageSock.SendMessage("playerFrame", f.FrameIdx, ":", p.PlayerNum, ":", p.Home)
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

	//Flush system by pushing <buffer size> + <filter size> frames through
	for range testFrameCount + storage.FrameBufferSize + 5 {
		frame = testutils.RandomNextFrame(frame)
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

	for _, f := range testFrames {
		topic, err := filterSocket.Recv(zmq.SNDMORE)
		assert.NoError(t, err)
		assert.Equal(t, "frame", topic)
		packet, err := filterSocket.RecvMessage(0)
		if err != nil {
			panic(err)
		}
		message := strings.Join(packet, "")
		assert.Greater(t, len(message), 0)
		filteredFrame := types.SmallFrame{}
		err = json.Unmarshal([]byte(message), &filteredFrame)
		assert.NoError(t, err)
		assert.Equal(t, f.FrameIdx, filteredFrame.FrameIdx, "Frames should arrive in correct order")
		for _, p := range f.Players {
			var found types.PlayerPosition
			for _, x := range filteredFrame.Players {
				if x.PlayerNum == p.PlayerNum && x.Home == p.Home {
					found = x
				}
			}

			assert.NotNil(t, found, "Player should exist in filtered frame")
			if filteredFrame.FrameIdx > 2 { //Skip the first filtered frames, since SavGol will not change those
				assert.NotEqual(t, p.X, found.X, "Filter should take effect on player X coordinate")
				assert.NotEqual(t, p.Y, found.Y, "Filter should take effect on player Y coordinate")
			}
		}
	}
}
