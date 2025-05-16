package testutils

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/google/uuid"
	"math/rand"
)

const fieldWidth, fieldHeight = 105, 68

func randomPosition() []float64 {
	return []float64{rand.Float64() * fieldWidth, rand.Float64() * fieldHeight, rand.Float64()}
}

func randomMove(pos []float64) {
	for i, value := range pos {
		pos[i] = value + rand.Float64() - 0.5
	}
}

func randomPlayer(number int) types.Player {
	return types.Player{
		Number:   fmt.Sprintf("%d", number),
		OptaId:   "",
		PlayerId: uuid.New().String(),
		Speed:    rand.Float64(),
		Xyz:      randomPosition(),
	}
}

func popRandom(nums []int) ([]int, int) {
	randomIndex := rand.Intn(len(nums))
	num := nums[randomIndex]
	newNums := make([]int, 0)
	newNums = append(newNums, nums[:randomIndex]...)
	newNums = append(newNums, nums[randomIndex+1:]...)
	return newNums, num
}

func RandomFrame(awayPlayers, homePlayers int) types.Frame {
	frame := types.Frame{
		AwayPlayers: make([]types.Player, awayPlayers),
		HomePlayers: make([]types.Player, homePlayers),
		Ball: struct {
			Speed float64   `json:"speed"`
			Xyz   []float64 `json:"xyz"`
		}{Speed: rand.Float64(), Xyz: randomPosition()},
		FrameIdx:  rand.Int(),
		GameClock: rand.Float64(),
		Period:    1,
		WallClock: int64(rand.Uint64()),
	}

	awayNumbers := make([]int, awayPlayers)
	homeNumbers := make([]int, homePlayers)

	for i := range awayNumbers {
		awayNumbers[i] = i
	}
	for i := range homeNumbers {
		homeNumbers[i] = i
	}

	for i := 0; i < awayPlayers; i++ {
		var num int
		awayNumbers, num = popRandom(awayNumbers)
		frame.AwayPlayers[i] = randomPlayer(num)
	}

	for i := 0; i < homePlayers; i++ {
		var num int
		homeNumbers, num = popRandom(homeNumbers)
		frame.HomePlayers[i] = randomPlayer(num)
	}

	return frame
}

func RandomNextFrame(previous types.Frame) types.Frame {
	next := types.Frame(previous)
	next.FrameIdx++
	for i := range next.AwayPlayers {
		randomMove(next.AwayPlayers[i].Xyz)
	}

	for i := range next.HomePlayers {
		randomMove(next.HomePlayers[i].Xyz)
	}

	randomMove(next.Ball.Xyz)
	next.WallClock++
	next.GameClock += rand.Float64()
	return next
}

func RandomFrameRange(teamSize int, frameCount int) []types.Frame {
	initFrame := RandomFrame(teamSize, teamSize)
	frames := make([]types.Frame, frameCount)
	for i := range frameCount {
		if i == 0 {
			frames[0] = initFrame
			continue
		}
		frames[i] = RandomNextFrame(frames[i-1])
	}
	return frames
}
