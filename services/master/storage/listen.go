package storage

import (
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strings"
	"sync"
)

func (w *Worker) listenConsume(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		topic, err := w.socketConsume.Recv(zmq.SNDMORE)
		if err != nil {
			w.Logger.Error("Failed to receive topic")
		}
		message, err := w.socketConsume.RecvMessage(0)
		if err != nil {
			w.Logger.Error("Error receiving message:", "error", err.Error())
		}

		if topic == "frame" {
			frame := types.DeserializeFrame(strings.Join(message, ""))
			for _, player := range frame.Players {
				w.mutex.Lock()
				buffer := w.players[player.PlayerId]
				buffer.Insert(player)
				w.mutex.Unlock()
			}
		}
	}
}

func (w *Worker) listenProvide(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		topic, err := w.socketProvide.Recv(zmq.SNDMORE)
		if err != nil {
			w.Logger.Error("Failed to read topic from message")
		}

		if topic == "playerFrame" {

		}

		if topic == "playerRange" {

		}

		if topic == "playerNumber" {

		}
	}
}
