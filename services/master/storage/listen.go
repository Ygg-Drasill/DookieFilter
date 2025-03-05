package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"strconv"
	"strings"
	"sync"
)

func (w *StorageWorker) listenConsume(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		message, err := w.socketConsume.RecvMessage(0)
		if err != nil {
			w.Logger.Error("Error receiving message:", "error", err.Error())
		}
		topic := message[0]
		if topic == "player" {
			frameIdx, err := strconv.Atoi(message[1])
			if err != nil {
				w.Logger.Error("Failed to parse player id", "error", err.Error())
				return
			}
			playerBuffer := w.players[message[2]]
			xy := strings.Split(message[3], ";")
			x, err := strconv.ParseFloat(xy[0], 64)
			y, err := strconv.ParseFloat(xy[1], 64)
			if err != nil {
				w.Logger.Error("Failed to parse player coordinates", "error", err.Error())
				return
			}
			position := types.PlayerPosition{
				FrameIdx: frameIdx,
				X:        x,
				Y:        y,
			}

			playerBuffer.Insert(position)
			w.Logger.Debug("Player buffer status", "count", playerBuffer.Count())
		}
	}
}

func (w *StorageWorker) listenProvide(wg *sync.WaitGroup) {
	defer wg.Done()
}
