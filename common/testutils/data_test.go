package testutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomMove(t *testing.T) {
	beforeMove := randomPosition()
	afterMove := make([]float64, len(beforeMove))
	copy(afterMove, beforeMove)
	assert.Equal(t, beforeMove, afterMove, "Expected before and after move to be same position without movement")
	randomMove(afterMove)
	assert.NotEqual(t, beforeMove, afterMove, "Expected position to change after movement")
}

func TestRandomPosition(t *testing.T) {
	a := randomPosition()
	b := randomPosition()
	assert.NotEqual(t, a, b, "Expected different positions from two different calls")
}

func TestRandomNextFrame(t *testing.T) {
	before := RandomFrame(3, 3)
	after := RandomNextFrame(before)
	assert.NotEqual(t, before, after, "Expected after frame to change from before frame")
}
