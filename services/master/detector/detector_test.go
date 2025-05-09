package detector

import (
	"sync"
	"testing"
	"time"

	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
)

func TestDetectHole(t *testing.T) {
    ctx, err := zmq.NewContext()
    socketInput, err := ctx.NewSocket(zmq.PUSH)
    socketImputation, err := ctx.NewSocket(zmq.PULL)
    err = socketInput.Connect(endpoints.InProcessEndpoint(endpoints.DETECTOR))
    assert.NoError(t, err)
    err = socketImputation.Bind("tcp://127.0.0.1:5555")
    assert.NoError(t, err)

    assert.NoError(t, err)
    w := New(ctx)
    wg := &sync.WaitGroup{}
    go w.Run(wg)

    //Generate random frame for testing
    frame := testutils.RandomFrame(2, 1)
    socketInput.SendMessage("frame", types.SerializeFrame(types.SmallFromBigFrame(frame)), 0)

    //Remove one player from the home list
    next := testutils.RandomNextFrame(frame)
    next.HomePlayers = next.HomePlayers[0: len(next.HomePlayers)-1]
    socketInput.SendMessage("frame", types.SerializeFrame(types.SmallFromBigFrame(next)), 0)


    doneChan := make(chan bool)
    go func() {
        message, err := socketImputation.RecvMessage(0)
        assert.NoError(t, err)
        assert.Greater(t, len(message), 0)
        doneChan <- true
    }()

    //Fail the test after 1 second
    go func() {
        time.Sleep(1 * time.Second)
        doneChan <- false
    }()

    result := <- doneChan
    assert.True(t, result, "expected messeage from detector worker within 1 second")
}
