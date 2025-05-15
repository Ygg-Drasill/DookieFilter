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

func randomPlayer() types.Player {
	number := fmt.Sprintf("%d", rand.Intn(99))
	return types.Player{
		Number:   number,
		OptaId:   "",
		PlayerId: uuid.New().String(),
		Speed:    rand.Float64(),
		Xyz:      randomPosition(),
	}
}

func RandomFrame(awayPlayers, homePlayers int) types.Frame {
	frame := types.Frame{
		AwayPlayers: make([]types.Player, 0),
		HomePlayers: make([]types.Player, 0),
		Ball: struct {
			Speed float64   `json:"speed"`
			Xyz   []float64 `json:"xyz"`
		}{Speed: rand.Float64(), Xyz: randomPosition()},
		FrameIdx:  rand.Int(),
		GameClock: rand.Float64(),
		Period:    1,
		WallClock: int64(rand.Uint64()),
	}

	for i := 0; i < awayPlayers; i++ {
		frame.AwayPlayers = append(frame.AwayPlayers, randomPlayer())
	}

	for i := 0; i < homePlayers; i++ {
		frame.HomePlayers = append(frame.HomePlayers, randomPlayer())
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
