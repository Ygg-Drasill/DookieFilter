package filter

import (
	"encoding/json"
	"errors"
	"github.com/Ygg-Drasill/DookieFilter/common/filter"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strings"
)

func (w *Worker) listenInput() {
	for {
		topic, err := w.socketInput.Recv(zmq.SNDMORE)
		if err != nil {
			w.Logger.Error("Failed to receive topic")
		}
		message, err := w.socketInput.RecvMessage(0)
		if err != nil {
			w.Logger.Error("Error receiving message:", "error", err.Error())
		}

		if topic == "frame" {
			frame := types.DeserializeFrame(strings.Join(message, ""))
			ff := filterableFrame(frame)
			filteredFrame, err := w.filter.Step(&ff)
			if (errors.Is(err, filter.NotFullError{})) {
				continue
			}

			byte, err := json.Marshal(filteredFrame)
			if err != nil {
				w.Logger.Error("Error marshalling filtered frame", "error", err.Error())
			}
			_, err = w.socketOutput.SendMessage(topic, byte)

			if err != nil {
				w.Logger.Error("Error sending message:", "error", err.Error())
			}
		}
	}
}
