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
	if frame.HomePlayers == nil ||
		frame.AwayPlayers == nil ||
		frame.Ball.Xyz == nil {
		return SmallFrame{}
	}
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
		playerData := fmt.Sprintf("%s;%g;%g", player.PlayerId, player.X, player.Y)
		parts[i] = playerData
	}

	ballString := fmt.Sprintf("%g;%g", frame.Ball.X, frame.Ball.Y)

	return fmt.Sprintf("%d:%s:%s", frameIdx, ballString, strings.Join(parts, ","))
}

func DeserializeFrame(data string) SmallFrame {
	frame := SmallFrame{
		Players: make([]PlayerPosition, 0),
		Ball:    Position{},
	}
	if data == "" {
		return frame
	}
	frameIdxAndData := strings.Split(data, ":")
	frameIdx, _ := strconv.Atoi(frameIdxAndData[0])
	frame.FrameIdx = frameIdx

	if frameIdxAndData[1] != "" {
		ballData := strings.Split(frameIdxAndData[1], ";")
		frame.Ball.X, _ = strconv.ParseFloat(ballData[0], 64)
		frame.Ball.Y, _ = strconv.ParseFloat(ballData[1], 64)
	}

	if frameIdxAndData[2] != "" {
		allPlayersData := strings.Split(frameIdxAndData[2], ",")
		frame.Players = make([]PlayerPosition, len(allPlayersData))
		for i, playerPart := range allPlayersData {
			playerData := strings.Split(playerPart, ";")
			x, _ := strconv.ParseFloat(playerData[1], 64)
			y, _ := strconv.ParseFloat(playerData[2], 64)
			player := PlayerPosition{
				FrameIdx: frameIdx,
				PlayerId: playerData[0],
				FrameIdx: frameIdx,
				Position: Position{
					X: x,
					Y: y,
				},
			}
			frame.Players[i] = player
		}
	}

	return frame
}
