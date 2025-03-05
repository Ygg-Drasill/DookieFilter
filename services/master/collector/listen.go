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
		frame := &types.Frame[types.DataPlayer]{}
		fullMessage := strings.Join(msg, "")
		err := json.Unmarshal([]byte(fullMessage), frame)
		if err != nil {
			return err
		}

		data := frame.Data[0]
		frameIdx := data.FrameIdx
		err = w.forwardPlayerPositions(frameIdx, data.HomePlayers)
		if err != nil {
			return err
		}
		err = w.forwardPlayerPositions(frameIdx, data.AwayPlayers)
		if err != nil {
			return err
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

func (w *CollectorWorker) forwardPlayerPositions(frameIdx int, players []types.Player) error {
	for _, p := range players {

		message := []any{
			"player",
			frameIdx,
			p.PlayerId,
			fmt.Sprintf("%f;%f", p.Xyz[0], p.Xyz[1]),
		}

		_, err := w.socketStore.SendMessage(message)
		if err != nil {
			return err
		}
		_, err = w.socketForward.SendMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}
