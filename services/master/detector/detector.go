package detector

import (
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"math"
	"strings"
	"sync"
)

type Worker struct {
	worker.BaseWorker

	socketListen *zmq.Socket

	stateBuffer *pringleBuffer.PringleBuffer[types.SmallFrame]
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
	w := &Worker{
		BaseWorker:  worker.NewBaseWorker(ctx, "detector"),
		stateBuffer: pringleBuffer.New[types.SmallFrame](10),
	}
	for _, opt := range options {
		opt(w)
	}

	return w
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer w.close()
	w.Logger.Info("Starting detector worker")
	err := w.connect()
	if err != nil {
		w.Logger.Error("Failed to bind/connect zmq sockets", "error", err)
	}

	for {
		topic, _ := w.socketListen.Recv(zmq.SNDMORE)
		if topic == "frame" {
			message, _ := w.socketListen.RecvMessage(0)
			frame := types.DeserializeFrame(strings.Join(message, ""))
			w.stateBuffer.Insert(frame)
			w.detect(frame)
			w.detectHoles(frame)
		}
	}
}

const (
	JumpThreshold = 5 //TODO: change me
	HoleThreshold = 10 // Number of consecutive frames a player can be missing before it's considered a hole
)

func (w *Worker) detect(frame types.SmallFrame) {
	prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - 1))
	if err != nil {
		w.Logger.Warn("No previous frame to compare")
		return
	}
	compareMap := make(map[string][]types.PlayerPosition)
	for _, player := range prevFrame.Players {
		_, ok := compareMap[player.PlayerId]
		if ok {
			compareMap[player.PlayerId][0] = player
		}
		if !ok {
			compareMap[player.PlayerId] = make([]types.PlayerPosition, 2)
			compareMap[player.PlayerId][0] = player
		}
	}
	for _, player := range frame.Players {
		_, ok := compareMap[player.PlayerId]
		if ok {
			compareMap[player.PlayerId][1] = player
		}
		if !ok {
			compareMap[player.PlayerId] = make([]types.PlayerPosition, 2)
			compareMap[player.PlayerId][1] = player
		}
	}

	for playerId, values := range compareMap {
		xDiff := math.Abs(values[0].Position.X - values[1].Position.X)
		yDiff := math.Abs(values[0].Position.Y - values[1].Position.Y)
		if xDiff > JumpThreshold || yDiff > JumpThreshold {
			w.Logger.Info("Jump detected", "player_id", playerId, "x_diff", xDiff, "y_diff", yDiff, "frame", frame.FrameIdx)
		}
	}
}

func (w *Worker) detectHoles(frame types.SmallFrame) {
	// Get previous frames from the buffer
	prevFrames := make([]types.SmallFrame, 0)
	for i := 1; i <= HoleThreshold; i++ {
		prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - i))
		if err != nil {
			w.Logger.Warn("No previous frame to compare")
			return
		}
	}

	if len(prevFrames) == 0 {
		return
	}

	// Track players that appear in current frame
	currentPlayers := make(map[string]bool)
	for _, player := range frame.Players {
		currentPlayers[player.PlayerId] = true
	}

	// Check each previous frame for missing players
	for _, prevFrame := range prevFrames {
		for _, player := range prevFrame.Players {
			if !currentPlayers[player.PlayerId] {
				w.Logger.Info("Player missing in frame", 
					"player_id", player.PlayerId,
					"frame", frame.FrameIdx,
					"last_seen_frame", prevFrame.FrameIdx)
			}
		}
	}
}
