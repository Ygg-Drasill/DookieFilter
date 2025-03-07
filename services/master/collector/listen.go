package collector

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strings"
)

func (w *Worker) listen() error {
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
		rawFrame := types.Frame{}
		fullMessage := strings.Join(msg, "")
		err := json.Unmarshal([]byte(fullMessage), &rawFrame)
		if err != nil {
			return err
		}

		frame := types.SmallFromBigFrame(rawFrame)

		err = w.forwardFrame(frame)
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

func (w *Worker) forwardFrame(frame types.SmallFrame) error {
	message := []any{
		"frame",
		types.SerializeFrame(frame),
	}

	_, err := w.socketStore.SendMessage(message...)
	if err != nil {
		return err
	}

	_, err = w.socketDetector.SendMessage(message...)
	if err != nil {
		return err
	}

	return nil
}
