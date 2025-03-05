package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"strings"
	"sync"
)

func (w *Worker) listenConsume(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		message, err := w.socketConsume.RecvMessage(0)
		if err != nil {
			w.Logger.Error("Error receiving message:", "error", err.Error())
		}
		topic := message[0]
		if topic == "frame" {
			frame := types.DeserializeFrame(strings.Join(message[1:], ""))
			for _, player := range frame.Players {
				buffer := w.players[player.PlayerId]
				buffer.Insert(player)
			}
		}
	}
}

func (w *Worker) listenProvide(wg *sync.WaitGroup) {
	defer wg.Done()
}
