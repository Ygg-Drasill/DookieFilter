package types

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"strconv"
	"strings"
)

type SmallFrame struct {
	FrameIdx int
	Players  []PlayerPosition
	Ball     Position
}

type Position struct {
	X float64
	Y float64
}

type PlayerPosition struct {
	Position
	FrameIdx int
	PlayerId string
}

func (pp PlayerPosition) Key() pringleBuffer.Key {
	return pringleBuffer.Key(pp.FrameIdx)
}

func (frame SmallFrame) Key() pringleBuffer.Key {
	return pringleBuffer.Key(frame.FrameIdx)
}

func SmallFromBigFrame(frame Frame) SmallFrame {
	smallFrame := SmallFrame{
		FrameIdx: frame.FrameIdx,
		Players:  make([]PlayerPosition, 0),
		Ball: Position{
			X: frame.Ball.Xyz[0],
			Y: frame.Ball.Xyz[1],
		},
	}

	for _, player := range frame.HomePlayers {
		smallFrame.Players = append(smallFrame.Players, PositionFromPlayer(player, frame.FrameIdx))
	}
	for _, player := range frame.AwayPlayers {
		smallFrame.Players = append(smallFrame.Players, PositionFromPlayer(player, frame.FrameIdx))
	}

	return smallFrame
}

func PositionFromPlayer(player Player, frameIdx int) PlayerPosition {
	return PlayerPosition{
		FrameIdx: frameIdx,
		PlayerId: player.PlayerId,
		Position: Position{
			X: player.Xyz[0],
			Y: player.Xyz[1],
		},
	}
}

func SerializeFrame(frame SmallFrame) string {
	frameIdx := frame.FrameIdx
	parts := make([]string, len(frame.Players))

	for i, player := range frame.Players {
		playerData := fmt.Sprintf("%s;%f;%f", player.PlayerId, player.X, player.Y)
		parts[i] = playerData
	}

	ballString := fmt.Sprintf("%f;%f", frame.Ball.X, frame.Ball.Y)

	return fmt.Sprintf("%d:%s:%s", frameIdx, ballString, strings.Join(parts, ","))
}

func DeserializeFrame(data string) SmallFrame {
	frameIdxAndData := strings.Split(data, ":")
	frameIdx, _ := strconv.Atoi(frameIdxAndData[0])
	allPlayerParts := strings.Split(frameIdxAndData[1], ",")
	players := make([]PlayerPosition, len(allPlayerParts))
	for i, playerPart := range allPlayerParts {
		playerData := strings.Split(playerPart, ";")
		x, _ := strconv.ParseFloat(playerData[1], 64)
		y, _ := strconv.ParseFloat(playerData[2], 64)
		player := PlayerPosition{
			PlayerId: playerData[0],
			Position: Position{
				X: x,
				Y: y,
			},
		}
		players[i] = player
	}

	return SmallFrame{
		FrameIdx: frameIdx,
		Players:  players,
		Ball:     Position{},
	}
}
