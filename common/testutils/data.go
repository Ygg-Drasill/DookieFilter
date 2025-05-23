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

	for range awayPlayers {
		frame.AwayPlayers = append(frame.AwayPlayers, randomPlayer())
	}

	for range homePlayers {
		frame.HomePlayers = append(frame.HomePlayers, randomPlayer())
	}

	return frame
}

func RandomNextFrame(previous types.Frame) types.Frame {
	next := types.Frame{}
	next.AwayPlayers = make([]types.Player, len(previous.AwayPlayers))
	next.HomePlayers = make([]types.Player, len(previous.HomePlayers))
	copy(next.AwayPlayers, previous.AwayPlayers)
	for i := range previous.AwayPlayers {
		next.AwayPlayers[i].Xyz = make([]float64, len(next.AwayPlayers[i].Xyz))
		copy(next.AwayPlayers[i].Xyz, previous.AwayPlayers[i].Xyz)
	}

	copy(next.HomePlayers, previous.HomePlayers)
	for i := range previous.HomePlayers {
		next.HomePlayers[i].Xyz = make([]float64, len(next.HomePlayers[i].Xyz))
		copy(next.HomePlayers[i].Xyz, previous.HomePlayers[i].Xyz)
	}

	next.FrameIdx++
	for i := range next.AwayPlayers {
		randomMove(next.AwayPlayers[i].Xyz)
	}

	for i := range next.HomePlayers {
		randomMove(next.HomePlayers[i].Xyz)
	}
	next.Ball.Xyz = make([]float64, len(previous.Ball.Xyz))
	copy(next.Ball.Xyz, previous.Ball.Xyz)
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
