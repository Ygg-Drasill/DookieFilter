package types

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"strconv"
	"strings"
)

type Player struct {
	Number   string    `json:"number"`
	OptaId   string    `json:"optaId"`
	PlayerId string    `json:"playerId"`
	Speed    float64   `json:"speed"`
	Xyz      []float64 `json:"xyz"`
}

type PlayerPosition struct {
	FrameIdx int     `json:"frameIdx"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

func (pp PlayerPosition) Key() pringleBuffer.Key {
	return pringleBuffer.Key(pp.FrameIdx)
}

func PositionFromPlayer(player Player, frameIdx int) PlayerPosition {
	return PlayerPosition{
		FrameIdx: frameIdx,
		X:        player.Xyz[0],
		Y:        player.Xyz[1],
	}
}

func SerializePlayerPositions(players []PlayerPosition) string {
	frameIdx := players[0].FrameIdx
	parts := make([]string, len(players))

	for i, player := range players {
		if player.FrameIdx != frameIdx {
			//TODO: something is big not good :(
		}
		playerData := fmt.Sprintf("%f;%f", player.X, player.Y)
		parts[i] = playerData
	}

	return fmt.Sprintf("%d:%s", frameIdx, strings.Join(parts, ","))
}

func DeserializePlayerPositions(data string) []PlayerPosition {
	frameIdxAndData := strings.Split(data, ":")
	frameIdx, _ := strconv.Atoi(frameIdxAndData[0])
	playerParts := strings.Split(frameIdxAndData[1], ",")
	positions := make([]PlayerPosition, len(playerParts))
	for i, playerPart := range playerParts {
		positionData := strings.Split(playerPart, ";")
		x, _ := strconv.ParseFloat(positionData[0], 64)
		y, _ := strconv.ParseFloat(positionData[1], 64)
		playerPosition := PlayerPosition{
			FrameIdx: frameIdx,
			X:        x,
			Y:        y,
		}
		positions[i] = playerPosition
	}
	return positions
}
