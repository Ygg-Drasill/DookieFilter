package storage

import (
	"fmt"
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
			w.Logger.Error("Failed to read topic from message", "error", err)
		}

		if topic == "playerFrame" {
			w.receivePlayerFrameRequest()
		}

		if topic == "playerNumber" {

		}

		if topic == "frameRange" {
			w.receiveFrameRangeRequest()
		}
	}
}

func (w *Worker) receivePlayerFrameRequest() {
	messageParts, err := w.socketProvide.RecvMessage(0)
	if err != nil {
		w.Logger.Error("Failed to receive request", "type", "playerFrame", "error", err)
		return
	}
	message := strings.Join(messageParts, "")
	frameIndexAndPlayerId := strings.Split(message, ":")
	frameIndex, playerId := frameIndexAndPlayerId[0], frameIndexAndPlayerId[1]
	w.Logger.Debug("Handling playerFrame request", "frameIndex", frameIndex, "playerId", playerId)
}

func (w *Worker) receiveFrameRangeRequest() {
	messageParts, err := w.socketProvide.RecvMessage(0)
	if err != nil {
		w.Logger.Error("Failed to receive request", "type", "frameRange", "error", err)
		return
	}
	message := strings.Join(messageParts, "")
	fmt.Println(message) //TODO: implement
}
