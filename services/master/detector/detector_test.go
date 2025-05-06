package detector_test

import (
	"testing"

	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/detector"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
)

func TestDetectHoles(t *testing.T) {
	// Step 1: Create a new detector Worker
	ctx, _ := zmq.NewContext()
	w := detector.New(ctx)
	w.HoleFlags = make(map[string]bool)
	w.HoleCount = 0

	// Initialize the stateBuffer
	w.StateBuffer = pringleBuffer.New[types.SmallFrame](10)

	// Step 2: Add mock frames to stateBuffer
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

	// Step 3: Call detectHoles with the current frame
	w.DetectHoles(currentFrame)

	// Step 4: Validate the outcome
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
