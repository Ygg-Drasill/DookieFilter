package detector_test

import (
	"encoding/json"
	"testing"

	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
)

func TestDetectHoles(t *testing.T) {
	//Create a new detector Worker
	ctx, _ := zmq.NewContext()
	w := detector.New(ctx)
	w.HoleFlags = make(map[string]bool)
	w.HoleCount = 0

	// Initialize the stateBuffer
	w.StateBuffer = pringleBuffer.New[types.SmallFrame](10)

	//Add mock frames to stateBuffer
	previousFrame := types.SmallFrame{
		FrameIdx: 1,
		Players: []types.PlayerPosition{
			{
				PlayerId: "Player1",
				Position: types.Position{X: 1.5, Y: 3.2},
			},
			{
				PlayerId: "Player2",
				Position: types.Position{X: 4.4, Y: 5.1},
			},
		},
	}

	currentFrame := types.SmallFrame{
		FrameIdx: 2,
		Players: []types.PlayerPosition{
			{
				PlayerId: "Player1",
				Position: types.Position{X: 2.0, Y: 3.8},
			},
			// Player2 is missing
		},
	}

	// Insert the previous frame into the buffer
	w.StateBuffer.Insert(previousFrame)

	//Call detectHoles with the current frame
	w.DetectHoles(currentFrame)

	//Validate the outcome
	// Assert that Player2 is marked as missing in holeFlags
	assert.True(t, w.HoleFlags["Player2"], "Player2 should be marked as missing in holeFlags")

	// Assert that Player1 is not marked as missing
	_, ok := w.HoleFlags["Player1"]
	assert.False(t, ok, "Player1 should not appear in holeFlags")

	// Assert that the holeCount has incremented
	assert.Equal(t, 1, w.HoleCount, "holeCount should be incremented by 1 for Player2")

	// Check that there is exactly one entry in holeFlags
	assert.Equal(t, 1, len(w.HoleFlags), "holeFlags should only have one missing player")
}

func TestSocketSend(t *testing.T) {
	// Step 1: Set up the Worker with ZeroMQ context
	ctx, _ := zmq.NewContext()
	defer func(ctx *zmq.Context) {
		err := ctx.Term()
		if err != nil {

		}
	}(ctx) // Ensure context is terminated after the test

	// Create a new worker
	w := detector.New(ctx)

	// Step 2: Create a ZeroMQ socket for testing
	receiver, _ := ctx.NewSocket(zmq.PULL) // Receiver socket
	defer func(receiver *zmq.Socket) {
		err := receiver.Close()
		if err != nil {

		}
	}(receiver)

	// Use an in-process endpoint for testing
	testEndpoint := "inproc://test"
	err := receiver.Bind(testEndpoint)
	assert.NoError(t, err, "Receiver should bind to the test endpoint")

	// Step 3: Set up the sender socket in the Worker
	w.SocketSend, err = w.SocketContext.NewSocket(zmq.PUSH)
	assert.NoError(t, err, "Should create PUSH socket")
	defer func(SocketSend *zmq.Socket) {
		err := SocketSend.Close()
		if err != nil {

		}
	}(w.SocketSend)

	err = w.SocketSend.Connect(testEndpoint)
	assert.NoError(t, err, "Sender should connect to the test endpoint")

	// Step 4: Prepare a test frame to send
	frame := types.SmallFrame{
		FrameIdx: 1,
		Players: []types.PlayerPosition{
			{PlayerId: "Player1", Position: types.Position{X: 1.4, Y: 2.7}},
		},
	}

	// Marshal the frame into JSON
	message, err := json.Marshal(frame)
	assert.NoError(t, err, "Should marshal frame to JSON")

	// Send the message
	_, err = w.SocketSend.SendMessage("frame", message)
	assert.NoError(t, err, "Should send the message without error")

	// Step 5: Receive and verify the message on the receiver socket
	receivedParts, err := receiver.RecvMessageBytes(0)
	assert.NoError(t, err, "Should receive a message without error")
	assert.Equal(t, 2, len(receivedParts), "Message should have 2 parts (topic + body)")

	// Verify the topic
	assert.Equal(t, "frame", string(receivedParts[0]), "Topic should be 'frame'")

	// Verify the message body
	var receivedFrame types.SmallFrame
	err = json.Unmarshal(receivedParts[1], &receivedFrame)
	assert.NoError(t, err, "Should unmarshal the received JSON")
	assert.Equal(t, frame, receivedFrame, "Sent and received frames should match")
}
