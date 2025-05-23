package testutils

import (
	"fmt"
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

	assert.Equal(t, before.FrameIdx+1, after.FrameIdx, "After frame should have same frame idx")
}

func TestPopRandom(t *testing.T) {
	count := 10
	nums := make([]int, count)
	for i := range nums {
		nums[i] = i
	}

	for range count {
		var num int
		nums, num = popRandom(nums)

		for i := range len(nums) {
			assert.NotEqual(t, nums[i], num)
		}
	}
}

func TestNoDuplicatePlayerNumbers(t *testing.T) {
	const count = 100
	for i := range 100 {
		t.Run(fmt.Sprintf("No Duplicate Player Numbers Run %d", i+1), func(t *testing.T) {
			frame := RandomFrame(count, count)
			homeNumbers := make(map[string]bool)
			awayNumbers := make(map[string]bool)
			for _, p := range frame.HomePlayers {
				_, ok := homeNumbers[p.Number]
				assert.False(t, ok, "Player number %s already exists in home", p.Number)
				homeNumbers[p.Number] = true
			}
			for _, p := range frame.AwayPlayers {
				_, ok := awayNumbers[p.Number]
				assert.False(t, ok, "Player number %s already exists in away", p.Number)
				awayNumbers[p.Number] = true
			}
			assert.Equal(t, count, len(homeNumbers))
			assert.Equal(t, count, len(awayNumbers))
		})
	}
}
