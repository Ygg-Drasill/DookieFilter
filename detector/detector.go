package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/types"
    "math"
    "strconv"
)

type Detector struct {
    Swaps     []Swap
    prevFrame frames
    playerMap map[int][]float64
}

type Swap struct {
    SwapID    int
    SwapFrame int
    p1        types.Player
    p2        types.Player
}

type player struct {
    playerID int
    xyz      []float64
}

type frames struct {
    players []player
    frameID int
}

func (d *Detector) Detect(frames []frames) {
    for _, f := range frames {
        d.DetectSwap(f)
    }
}

func (d *Detector) DetectSwap(currentFrame frames) {
    if d.prevFrame.frameID == 0 {
        d.prevFrame = currentFrame
        d.playerMap = make(map[int][]float64)
        for _, p := range currentFrame.players {
            d.playerMap[p.playerID] = p.xyz
        }
        return
    }

    for _, p1 := range currentFrame.players {
        prevP1, ok1 := d.playerMap[p1.playerID]
        if !ok1 {
            continue
        }

        for _, p2 := range currentFrame.players {
            if p1.playerID == p2.playerID {
                continue
            }

            prevP2, ok2 := d.playerMap[p2.playerID]
            if !ok2 {
                continue
            }

            if adjacencyThreshold(prevP1, prevP2) && adjacencyThreshold(p1.xyz, p2.xyz) {
                swap := Swap{
                    SwapID:    len(d.Swaps) + 1,
                    SwapFrame: currentFrame.frameID,
                    p1:        types.Player{PlayerId: strconv.Itoa(p1.playerID)},
                    p2:        types.Player{PlayerId: strconv.Itoa(p2.playerID)},
                }
                d.Swaps = append(d.Swaps, swap)
            }
        }
    }

    d.playerMap = make(map[int][]float64)
    for _, p := range currentFrame.players {
        d.playerMap[p.playerID] = p.xyz
    }
    d.prevFrame = currentFrame
}

func adjacencyThreshold(p1 []float64, p2 []float64) bool {
    adjT := 1.1
    return math.Abs(p1[0]-p2[0]) < adjT && math.Abs(p1[1]-p2[1]) < adjT
}
