package main

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
	"github.com/Ygg-Drasill/DookieFilter/services/master/storage"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testFrameCount = 10

func TestMasterServiceIntegration(t *testing.T) {
	//slog.SetLogLoggerLevel(slog.LevelError)
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

	initialFrame := testutils.RandomFrame(3, 3)
	for i, frame := 0, initialFrame; i < testFrameCount; i, frame = i+1, testutils.RandomNextFrame(frame) {
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

	time.Sleep(3 * time.Second)
}
