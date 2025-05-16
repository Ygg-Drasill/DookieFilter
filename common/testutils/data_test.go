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
