package types

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"log/slog"
	"strconv"
	"strings"
)

type SmallFrame struct {
	FrameIdx int              `json:"frameIdx"`
	Players  []PlayerPosition `json:"players"`
	Ball     Position         `json:"ball"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type PlayerKey struct {
	PlayerNumber int  `json:"playerNumber"`
	Home         bool `json:"home"`
}

func NewPlayerKey(num int, home bool) PlayerKey {
	return PlayerKey{num, home}
}

type PlayerPosition struct {
	Position
	FrameIdx  int  `json:"frameIdx"`
	PlayerNum int  `json:"number"`
	Home      bool `json:"home"`
}

func (pp PlayerPosition) Key() pringleBuffer.Key {
	return pringleBuffer.Key(pp.FrameIdx)
}

func (pp PlayerPosition) SKey() string {
	return fmt.Sprintf("%d:%t", pp.PlayerNum, pp.Home)
}

func DeSKey(k string) (int, bool) {
	x := strings.Split(k, ":")
	n, err := strconv.Atoi(x[0])
	if err != nil {
		return 0, false
	}
	h, err := strconv.ParseBool(x[1])
	return n, h
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
		smallFrame.Players = append(smallFrame.Players, PositionFromPlayer(player, frame.FrameIdx, true))
	}
	for _, player := range frame.AwayPlayers {
		smallFrame.Players = append(smallFrame.Players, PositionFromPlayer(player, frame.FrameIdx, false))
	}

	return smallFrame
}

func PositionFromPlayer(player Player, frameIdx int, home bool) PlayerPosition {
	number, err := strconv.Atoi(player.Number)
	if err != nil {
		slog.Error("Failed to convert player number to int", "error", err)
		return PlayerPosition{}
	}
	return PlayerPosition{
		FrameIdx:  frameIdx,
		PlayerNum: number,
		Home:      home,
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
		playerData := fmt.Sprintf("%d;%t;%g;%g",
			player.PlayerNum,
			player.Home,
			player.X,
			player.Y)
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
			x, _ := strconv.ParseFloat(playerData[2], 64)
			y, _ := strconv.ParseFloat(playerData[3], 64)
			playerNumber, _ := strconv.Atoi(playerData[0])
			home, _ := strconv.ParseBool(playerData[1])
			player := PlayerPosition{
				FrameIdx:  frameIdx,
				PlayerNum: playerNumber,
				Home:      home,
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
