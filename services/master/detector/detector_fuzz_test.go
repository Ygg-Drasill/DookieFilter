package detector

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
)

func FuzzSwap(f *testing.F) {
	for range 10 {
		f.Add(rand.Int63())
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))
		playerCount := r.Intn(1000) + 1

		playerMap := randomPlayerMap(r, playerCount)

		mockWorker := getMockWorker()
		swappers := mockWorker.swap(playerMap)

		for key := range swappers {
			assert.True(t, strings.HasSuffix(key, ":1"))
			_, ok := playerMap[key]
			assert.True(t, ok)
		}
		var swapCount, jumpCount int
		for _, v := range swappers {
			if v {
				swapCount++
			} else {
				jumpCount++
			}
		}

		assert.GreaterOrEqual(t, playerCount, swapCount)

		// ensure amount of swaps is even
		assert.Equalf(t, 0, swapCount%2, "Expected even amount of swaps, got %d", swapCount)

		assert.GreaterOrEqual(t, playerCount, jumpCount)
		assert.Equal(t, playerCount, swapCount+jumpCount)
		t.Logf("\nPlayer count: %d\nTotal swaps: %d\nTotal jumps: %d", playerCount, swapCount, jumpCount)
	})
}

func randomPlayerMap(r *rand.Rand, playerCount int) map[string]types.PlayerPosition {
	playerMap := make(map[string]types.PlayerPosition)
	for i := range playerCount {
		// prev frame
		playerId := fmt.Sprintf("player%d", i+1)
		x0, y0 := r.Float64()*10-5, r.Float64()*10-5
		key0 := fmt.Sprintf("%s:10:0", playerId)
		playerMap[key0] = types.PlayerPosition{
			PlayerId: playerId,
			FrameIdx: 10,
			Position: types.Position{X: x0, Y: y0},
		}
		// curr frame
		x1, y1 := r.Float64()*10-5, r.Float64()*10-5
		key1 := fmt.Sprintf("%s:11:1", playerId)
		playerMap[key1] = types.PlayerPosition{
			PlayerId: playerId,
			FrameIdx: 11,
			Position: types.Position{X: x1, Y: y1},
		}
	}
	return playerMap
}
