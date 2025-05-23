package testutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomMove(t *testing.T) {
	beforeMove := randomPosition()
	afterMove := make([]float64, len(beforeMove))
	copy(afterMove, beforeMove)
	for i := range afterMove {
		assert.Equal(t, beforeMove[i], afterMove[i], "Expected before and after move to be same position without movement")
	}
	randomMove(afterMove)
	for i := range afterMove {
		assert.NotEqual(t, beforeMove[i], afterMove[i], "Expected position to change after movement")
	}
}

func TestRandomPosition(t *testing.T) {
	a := randomPosition()
	b := randomPosition()
	assert.NotEqual(t, a, b, "Expected different positions from two different calls")
}

func TestRandomNextFrame(t *testing.T) {
	pCount := 3
	before := RandomFrame(pCount, pCount)
	after := RandomNextFrame(before)
	for i := range pCount {
		homePosBefore := before.HomePlayers[i].Xyz
		homePosAfter := after.HomePlayers[i].Xyz
		assert.NotEqual(t, homePosBefore[0], homePosAfter[0], "Expected after frame to change from before frame")
		assert.NotEqual(t, homePosBefore[1], homePosAfter[1], "Expected after frame to change from before frame")

		awayPosBefore := before.AwayPlayers[i].Xyz
		awayPosAfter := after.AwayPlayers[i].Xyz
		assert.NotEqual(t, awayPosBefore[0], awayPosAfter[0], "Expected after frame to change from before frame")
		assert.NotEqual(t, awayPosBefore[1], awayPosAfter[1], "Expected after frame to change from before frame")
	}
}
