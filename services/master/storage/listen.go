package storage

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"strings"
	"sync"
)

const (
	REQ_PLAYER_FRAME = "playerFrame"
	REQ_FRAME_RANGE  = "frameRange"
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
				buffer, ok := w.players[player.PlayerId]
				if !ok {
					playerBuffer := *pringleBuffer.New[types.PlayerPosition](w.bufferSize)
					w.players[player.PlayerId] = playerBuffer
					buffer = playerBuffer
				}

				buffer.Insert(player)
				w.players[player.PlayerId] = buffer
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
			w.handlePlayerFrameRequest()
		}

		if topic == "playerNumber" {

		}

		if topic == "frameRange" {
			w.handleFrameRangeRequest()
		}
	}
}

func (w *Worker) handlePlayerFrameRequest() {
	messageParts, err := w.socketProvide.RecvMessage(0)
	if err != nil {
		w.Logger.Error("Failed to receive request", "type", REQ_PLAYER_FRAME, "error", err)
		return
	}
	message := strings.Join(messageParts, "")
	frameIndexAndPlayerId := strings.Split(message, ":")
	frameIndexRaw, playerId := frameIndexAndPlayerId[0], frameIndexAndPlayerId[1]
	w.Logger.Debug("Handling storage request", "type", REQ_PLAYER_FRAME, "frameIndex", frameIndexRaw, "playerId", playerId)

	frameIndex, err := strconv.Atoi(frameIndexRaw)
	if err != nil {
		w.Logger.Error("Failed to convert frame index to int", "error", err)
		w.respondEmpty()
		return
	}
	playerBuffer := w.players[playerId]
	position, err := playerBuffer.Get(pringleBuffer.Key(frameIndex))
	if err != nil {
		w.Logger.Error("Failed to get position", "frameIndex", frameIndex, "playerId", playerId, "error", err)
		w.respondEmpty()
		return
	}

	response := fmt.Sprintf("%g;%g", position.X, position.Y)

	_, err = w.socketProvide.Send(response, 0)
	if err != nil {
		w.Logger.Error("Failed to send message", "type", REQ_PLAYER_FRAME, "error", err)
		w.respondEmpty()
		return
	}
}

func (w *Worker) handleFrameRangeRequest() {
	messageParts, err := w.socketProvide.RecvMessage(0)
	if err != nil {
		w.Logger.Error("Failed to receive request", "type", REQ_FRAME_RANGE, "error", err)
		return
	}
	message := strings.Join(messageParts, "")
	fmt.Println(message) //TODO: implement
}

func (w *Worker) respondEmpty() {
	w.socketProvide.Send("", 0)
}
