package storage

import (
	"errors"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

func (w *Worker) forwardFramesToFilter(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		correctedPlayer := <-w.correctedPlayersChan
		frame, err := w.frameBuffer.Get(pringleBuffer.Key(correctedPlayer.FrameIdx))
		if errors.Is(err, pringleBuffer.NotFoundError{}) || frame == nil {
			newFrame := &types.SmallFrame{
				FrameIdx: correctedPlayer.FrameIdx,
				Players:  make([]types.PlayerPosition, 1),
				Ball:     <-w.ballChan, //TODO: this is very flaky
			}
			newFrame.Players[0] = correctedPlayer
			w.frameBuffer.Insert(newFrame)
			continue
		}

		frame.Players = append(frame.Players, correctedPlayer)
	}
}

func (w *Worker) forwardToFilter() func(frame *types.SmallFrame) {
	return func(frame *types.SmallFrame) {
		_, err := w.socketFilter.Send("frame", zmq.SNDMORE)
		_, err = w.socketFilter.SendMessage(types.SerializeFrame(*frame))
		if err != nil {
			w.Logger.Error("Error sending message frame to filter", "error", err.Error())
		}
	}
}
