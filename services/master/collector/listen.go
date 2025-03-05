package collector

import (
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strings"
)

func (w *CollectorWorker) listen() error {
	topic, err := w.socketListen.Recv(zmq.SNDMORE)
	if err != nil {
		return err
	}

	msg, err := w.socketListen.RecvMessage(0)
	if err != nil {
		return err
	}

	if topic == "frame" {
		//parse raw frame
		frame := &types.GamePacket[types.Frame]{}
		fullMessage := strings.Join(msg, "")
		err := json.Unmarshal([]byte(fullMessage), frame)
		if err != nil {
			return err
		}

		for _, data := range frame.Data {
			frameIdx := data.FrameIdx
			allPlayers := make([]types.PlayerPosition, 0)
			allPlayers, err = w.storePlayerPositions(frameIdx, data.HomePlayers, allPlayers)
			if err != nil {
				return err
			}
			allPlayers, err = w.storePlayerPositions(frameIdx, data.AwayPlayers, allPlayers)
			if err != nil {
				return err
			}

			err = w.forwardPositionsToDetector(frameIdx, allPlayers)
			if err != nil {
				return err
			}
		}
	}

	if topic == "player" {
		_, err := w.socketStore.SendMessage(0, "player", msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *CollectorWorker) storePlayerPositions(frameIdx int, players []types.Player, allPlayers []types.PlayerPosition) ([]types.PlayerPosition, error) {
	for _, p := range players {

		position := types.PositionFromPlayer(p, frameIdx)
		message := []any{
			"player",
			frameIdx,
			p.PlayerId,
			fmt.Sprintf("%f;%f", position.X, position.Y),
		}

		_, err := w.socketStore.SendMessage(message)
		if err != nil {
			return allPlayers, err
		}
		allPlayers = append(allPlayers, position)
	}

	return allPlayers, nil
}

func (w *CollectorWorker) forwardPositionsToDetector(frameIdx int, players []types.PlayerPosition) error {
	message := []any{
		"playerframe",
		types.SerializePlayerPositions(players),
	}

	_, err := w.socketDetector.SendMessage(message)
	if err != nil {
		return err
	}

	return nil
}
