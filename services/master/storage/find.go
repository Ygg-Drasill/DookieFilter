package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"math"
	"sort"
)

func (w *Worker) findNearestPlayers(frameIndex int, count int, target types.PlayerKey) FrameNearest {
	bufferKey := pringleBuffer.Key(frameIndex)
	p := w.players[target]
	targetPos, err := p.Get(bufferKey)
	if err != nil {
		w.Logger.Error("Failed to find target player in buffer", "frameIndex", frameIndex, "error", err)
	}

	homeNearest := make(playerDistanceSorted, 0)
	awayNearest := make(playerDistanceSorted, 0)

	for key, buf := range w.players {
		if key == target {
			continue
		}

		position, err := buf.Get(bufferKey)
		if err != nil {
			w.Logger.Error("Failed to find player in buffer", "frameIndex", frameIndex, "error", err)
			continue
		}

		distanceX := math.Abs(position.X - targetPos.X)
		distanceY := math.Abs(position.Y - targetPos.Y)
		distance := distanceX + distanceY
		pd := playerDistance{
			PlayerPosition: position,
			distance:       distance,
		}
		if position.Home {
			homeNearest = append(homeNearest, pd)
		} else {
			awayNearest = append(awayNearest, pd)
		}
	}

	sort.Sort(homeNearest)
	sort.Sort(awayNearest)
	return FrameNearest{
		Target: targetPos,
		Home:   homeNearest[:count],
		Away:   awayNearest[:count],
	}
}
