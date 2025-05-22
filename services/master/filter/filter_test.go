package filter

import (
	"encoding/json"
	"math"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
)

func addNoise(v float64, scale float64) float64 {
	return v + (rand.Float64()-.5)*scale
}

func sin(i int) float64 {
	return math.Sin(float64(i)/framesCount) * 10
}

const framesCount = 10

func Test_FilterWorker(t *testing.T) {
	ctx, err := zmq.NewContext()
	assert.NoError(t, err)
	inputSocket, err := ctx.NewSocket(zmq.PUSH)
	assert.NoError(t, err)
	outputSocket, err := ctx.NewSocket(zmq.PULL)
	assert.NoError(t, err)
	endpoint := "inproc://test"

	worker := New(ctx, WithOutputEndpoint(endpoint))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go worker.Run(wg)

	inputSocket.Connect(endpoints.InProcessEndpoint(endpoints.FILTER_INPUT))
	outputSocket.Connect(endpoint)

	frames := make([]types.Frame, framesCount)
	frames[0] = testutils.RandomFrame(3, 3)
	//pNum := frames[0].HomePlayers[0].Number
	frames[0].HomePlayers[0].Xyz = []float64{
		addNoise(sin(0), 2),
		addNoise(sin(0), 2)}
	for i := 1; i < framesCount; i++ {
		frames[i] = testutils.RandomNextFrame(frames[i-1])
		frames[i].HomePlayers[0].Xyz = []float64{
			addNoise(sin(i), 2),
			addNoise(sin(i), 2)}
	}

	smallFrames := make([]types.SmallFrame, framesCount)
	for i, f := range frames {
		smallFrames[i] = types.SmallFromBigFrame(f)
	}

	for _, frame := range smallFrames {
		msg := types.SerializeFrame(frame)
		inputSocket.SendMessage("frame", msg)
	}

	filteredFrames := make([]types.SmallFrame, framesCount)
	for i := range framesCount {
		msg, _ := outputSocket.RecvMessage(0)
		frame := types.SmallFrame{}
		err = json.Unmarshal([]byte(strings.Join(msg[1:], "")), &frame)
		filteredFrames[i] = frame
	}

	rawErrorX, filterErrorX := .0, .0
	rawErrorY, filterErrorY := .0, .0

	for i := range framesCount - 1 {
		assert.Equal(t, smallFrames[i].Players[0].PlayerNum, filteredFrames[i].Players[0].PlayerNum, "Player order should not change")

		cleanSignal := sin(i)
		rawErrorX += math.Abs(cleanSignal - smallFrames[i].Players[0].X)
		rawErrorY += math.Abs(cleanSignal - smallFrames[i].Players[0].Y)
		filterErrorX += math.Abs(cleanSignal - filteredFrames[i].Players[0].X)
		filterErrorY += math.Abs(cleanSignal - filteredFrames[i].Players[0].X)
	}

	assert.Less(t, filterErrorX, rawErrorX, "Error from X should improve after filter")
	assert.Less(t, filterErrorY, rawErrorY, "Error from Y should improve after filter")
}
