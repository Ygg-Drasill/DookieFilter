package filter

import (
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"strings"
	"sync"
)

func (w *Worker) listenInput(wg *sync.WaitGroup) {
	defer wg.Done()
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
			for _, player := range frame.Players {
				//w.mutex.Lock()
			}

		}

	}
}
