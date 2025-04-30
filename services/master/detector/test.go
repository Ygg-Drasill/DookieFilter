package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"log"
)

// MockWorker is a simplified version of the detector worker for testing
type MockWorker struct {
	stateBuffer *pringleBuffer.PringleBuffer[types.SmallFrame]
}

func NewMockWorker() *MockWorker {
	return &MockWorker{
		stateBuffer: pringleBuffer.New[types.SmallFrame](10),
	}
}

func (w *MockWorker) detectHoles(frame types.SmallFrame) {
	// Get previous frames from the buffer
	prevFrames := make([]types.SmallFrame, 0)
	for i := 1; i <= 10; i++ {
		prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - i))
		if err == nil {
			prevFrames = append(prevFrames, prevFrame)
		}
	}

	if len(prevFrames) == 0 {
		return
	}

	// Track players that appear in current frame
	currentPlayers := make(map[string]bool)
	for _, player := range frame.Players {
		currentPlayers[player.PlayerId] = true
	}

	// Check each previous frame for missing players
	for _, prevFrame := range prevFrames {
		for _, player := range prevFrame.Players {
			if !currentPlayers[player.PlayerId] {
				log.Printf("Hole detected! Player %s missing in frame %d (last seen in frame %d)",
					player.PlayerId, frame.FrameIdx, prevFrame.FrameIdx)
			}
		}
	}
}

func main() {
	worker := NewMockWorker()
	
	// Create test frames with a hole
	frames := []types.SmallFrame{
		{
			FrameIdx: 1,
			Players: []types.PlayerPosition{
				{PlayerId: "player1", Position: types.Position{X: 0, Y: 0}},
				{PlayerId: "player2", Position: types.Position{X: 1, Y: 1}},
			},
		},
		{
			FrameIdx: 2,
			Players: []types.PlayerPosition{
				{PlayerId: "player1", Position: types.Position{X: 0, Y: 0}},
				// player2 is missing in this frame
			},
		},
		{
			FrameIdx: 3,
			Players: []types.PlayerPosition{
				{PlayerId: "player1", Position: types.Position{X: 0, Y: 0}},
				{PlayerId: "player2", Position: types.Position{X: 1, Y: 1}},
			},
		},
	}
	
	// Process frames
	for _, frame := range frames {
		worker.stateBuffer.Insert(frame)
		worker.detectHoles(frame)
	}
} 