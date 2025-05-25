package storage

import (
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"strings"
	"sync"
)

const (
	REQ_PLAYER_FRAME        = "playerFrame"
	REQ_FRAME_RANGE_NEAREST = "frameRangeNearest"
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

		if topic == "position" {
			player := types.PlayerPosition{}
			err = json.Unmarshal([]byte(strings.Join(message, "")), &player)
			if err != nil {
				w.Logger.Error("Error unmarshalling", "error", err.Error())
			}
			w.Logger.Info("Received new position", "player_number", player.PlayerNum, "player_home", player.Home)
			key := types.NewPlayerKey(player.PlayerNum, player.Home)
			w.mutex.Lock()
			buffer := w.players[key]
			buffer.Insert(player)
			w.players[key] = buffer
			w.mutex.Unlock()
		}

		if topic == "frame" {
			frame := types.DeserializeFrame(strings.Join(message, ""))
			//w.ballChan <- frame.Ball //TODO: fix
			for _, player := range frame.Players {
				w.mutex.Lock()
				key := types.NewPlayerKey(player.PlayerNum, player.Home)
				buffer, ok := w.players[key]
				if !ok {
					playerBuffer := *pringleBuffer.New[types.PlayerPosition](
						w.bufferSize,
						pringleBuffer.WithOnPopTail(func(element types.PlayerPosition) {
							w.correctedPlayersChan <- element
						}))
					w.players[key] = playerBuffer
					buffer = playerBuffer
				}

				buffer.Insert(player)
				w.players[key] = buffer
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
	}
}

func (w *Worker) listenAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		data, err := w.socketAPI.RecvMessage(0)
		topic := data[0]
		if err != nil {
			w.Logger.Error("Failed to read topic from message", "error", err, "topic", topic)
			return
		}
		message := strings.Join(data[1:], "")
		if topic == REQ_FRAME_RANGE_NEAREST {
			w.handleFrameRangeNearestRequest(message)
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
	frameIndexRaw, playerNumRaw, homeRaw := frameIndexAndPlayerId[0], frameIndexAndPlayerId[1], frameIndexAndPlayerId[2]
	w.Logger.Debug("Handling storage request", "type", REQ_PLAYER_FRAME, "frameIndex", frameIndexRaw, "playerNum", playerNumRaw, "home", homeRaw)

	frameIndex, err := strconv.Atoi(frameIndexRaw)
	playerNum, err := strconv.Atoi(playerNumRaw)
	home, err := strconv.ParseBool(homeRaw)

	if err != nil {
		w.Logger.Error("Failed to convert frame index to int", "error", err)
		w.respondEmpty()
		return
	}
	w.mutex.Lock()
	playerBuffer := w.players[types.NewPlayerKey(playerNum, home)]
	position, err := playerBuffer.Get(pringleBuffer.Key(frameIndex))
	w.mutex.Unlock()
	if err != nil {
		w.Logger.Error("Failed to get position", "frameIndex", frameIndex, "playerNum", playerNum, "home", home, "error", err)
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

func (w *Worker) handleFrameRangeNearestRequest(message string) {
	request := FrameRangeNearestRequest{}
	err := json.Unmarshal([]byte(message), &request)
	if err != nil {
		w.Logger.Error("Failed to unmarshal request", "error", err)
	}
	response := make(FrameRangeNearestResponse, request.EndIndex-request.StartIndex+1)
	for i := request.StartIndex; i <= request.EndIndex; i++ {
		response[i-request.StartIndex] = w.findNearestPlayers(i, request.NearestCount, request.TargetPlayer)
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		w.Logger.Error("Failed to marshal response", "error", err)
	}

	_, err = w.socketAPI.SendMessage(responseData)
}

func (w *Worker) respondEmpty() {
	_, err := w.socketProvide.Send("", 0)
	if err != nil {
		w.Logger.Error("Failed to send empty message", "error", err)
	}
}
