package detector

import (
	"encoding/json"
	"math"
	"strings"
	"sync"
	"reflect"
	"time"
	"testing"

	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
)

type Worker struct {
	worker.BaseWorker

	socketListen *zmq.Socket
	socketSend   *zmq.Socket

	stateBuffer *pringleBuffer.PringleBuffer[types.SmallFrame]
	holeFlags   map[string]bool
	holeCount   int
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
	prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - 1))
	if err != nil {
		w.Logger.Warn("No previous frame to compare")
		return
	}
	// Create sets for efficient lookup
	currentPlayers := make(map[string]bool)
	for _, player := range frame.Players {
		currentPlayers[player.PlayerId] = true
	}

	prevPlayers := make(map[string]bool)
	for _, player := range prevFrame.Players {
		prevPlayers[player.PlayerId] = true
	}

	// Check for players who were present before but are missing now
	for playerId := range prevPlayers {
		if !currentPlayers[playerId] {
			// Player is missing in the current frame
			if !w.holeFlags[playerId] {
				// Player just went missing, set the flag.
				w.holeFlags[playerId] = true // Set holeFlag to true when position is missing
				w.Logger.Info("HoleFlag: Player %s started missing at frame %d", "player_id", playerId, "frame", frame.FrameIdx)
				w.holeCount++ // Increment hole count when a player returns

			}
		}
	}
	var err error
	w.socketSend, err = w.SocketContext.NewSocket(zmq.PUSH)
	if err != nil {
		w.Logger.Error("Failed to create socket", "error", err)
		return
	}
	err = w.socketSend.Connect(endpoints.InProcessEndpoint(endpoints.COLLECTOR))
	if err != nil {
		w.Logger.Error("Failed to connect socket", "error", err)
		return
	}

	// Declare message first, then assign to existing err
	var message []byte
	message, err = json.Marshal(frame)
	if err != nil {
		w.Logger.Error("Failed to marshal frame to JSON", "error", err)
		return
	}

	// Declare messageLength first, then assign to existing err
	var messageLength int
	messageLength, err = w.socketSend.SendMessage("frame", message)
	if err != nil {
		// Use messageLength in the error log
		w.Logger.Error("Failed to send message", "length", messageLength, "error", err)
	}
}


