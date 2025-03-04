package detec

import (
	"errors"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"io"
	"math"
	"strconv"
)

type Detector struct {
	Swaps     []Swap
	prevFrame frames
	PlayerMap map[int][]float64
}

type Swap struct {
	SwapID    int
	SwapFrame int
	P1        types.Player
	P2        types.Player
}

type player struct {
	playerNumber int
	xyz          []float64
}

type frames struct {
	players []player
	frameID int
}

func (d *Detector) Detect(frames []frames) []Swap {
	for _, f := range frames {
		d.DetectSwap(f)
	}
	return d.Swaps
}

func (d *Detector) DetectSwap(currentFrame frames) {
	if d.prevFrame.frameID == 0 {
		d.prevFrame = currentFrame
		if currentFrame.players == nil {
			return
		}
		for _, p := range currentFrame.players {
			d.PlayerMap[p.playerNumber] = p.xyz
		}
		return
	}

	for _, p1 := range currentFrame.players {
		prevP1, ok1 := d.PlayerMap[p1.playerNumber]
		if !ok1 {
			continue
		}

		for _, p2 := range currentFrame.players {
			if p1.playerNumber == p2.playerNumber {
				continue
			}

			prevP2, ok2 := d.PlayerMap[p2.playerNumber]
			if !ok2 {
				continue
			}

			if adjacencyThreshold(prevP1, prevP2) && adjacencyThreshold(p1.xyz, p2.xyz) {
				swap := Swap{
					SwapID:    len(d.Swaps) + 1,
					SwapFrame: currentFrame.frameID,
					P1:        types.Player{Number: strconv.Itoa(p1.playerNumber)},
					P2:        types.Player{Number: strconv.Itoa(p2.playerNumber)},
				}
				d.Swaps = append(d.Swaps, swap)
			}
		}
	}

	d.PlayerMap = make(map[int][]float64)
	for _, p := range currentFrame.players {
		d.PlayerMap[p.playerNumber] = p.xyz
	}
	d.prevFrame = currentFrame
}

func adjacencyThreshold(p1 []float64, p2 []float64) bool {
	adjT := 1.1
	return math.Abs(p1[0]-p2[0]) < adjT && math.Abs(p1[1]-p2[1]) < adjT
}

func LoadFrames() ([]frames, error) {
	var allFrames []frames
	f, err := frameReader.New("test.jsonl")
	if err != nil {
		return nil, err
	}
	for {
		frame, err := f.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		frameData := frames{
			frameID: frame.Data[0].FrameIdx,
		}
		for _, p := range frame.Data[0].HomePlayers {
			frameData.players = append(frameData.players, player{
				playerNumber: parsePlayerID(p.Number),
				xyz:          p.Xyz,
			})
		}
		for _, p := range frame.Data[0].AwayPlayers {
			frameData.players = append(frameData.players, player{
				playerNumber: parsePlayerID(p.Number),
				xyz:          p.Xyz,
			})
		}
		allFrames = append(allFrames, frameData)
	}
	return allFrames, nil
}

func parsePlayerID(playerID string) int {
	id, err := strconv.Atoi(playerID)
	if err != nil {
		return -1
	}
	return id
}
