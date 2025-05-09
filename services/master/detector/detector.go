package detector

import (
    "encoding/json"
    "errors"
    "fmt"
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

    socketListen     *zmq.Socket
    socketImputation *zmq.Socket

    StateBuffer *pringleBuffer.PringleBuffer[types.SmallFrame]
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
    w := &Worker{
        BaseWorker:  worker.NewBaseWorker(ctx, "detector"),
        StateBuffer: pringleBuffer.New[types.SmallFrame](10),
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
            w.StateBuffer.Insert(frame)
            w.detect(frame)
            w.detectHoles(frame)
        }
    }
}

const (
    JumpThreshold = 5 //TODO: change me
)

func (w *Worker) detect(frame types.SmallFrame) {

    var checkPlayer = make(map[string]types.PlayerPosition)
    prevFrame, err := w.StateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - 1))
    if errors.Is(err, pringleBuffer.PringleBufferError{}) {
        w.Logger.Warn("No previous frame to compare")
        return
    }
    if err != nil {
        w.Logger.Error("Failed to get previous frame", "error", err)
        return
    }

    compareMap := make(map[string][]types.PlayerPosition)
    for _, player := range prevFrame.Players {
        _, ok := compareMap[player.PlayerId]
        if ok {
            player.FrameIdx = prevFrame.FrameIdx
            compareMap[player.PlayerId][0] = player
        }
        if !ok {
            compareMap[player.PlayerId] = make([]types.PlayerPosition, 2)
            player.FrameIdx = prevFrame.FrameIdx
            compareMap[player.PlayerId][0] = player
        }
    }
    for _, player := range frame.Players {
        _, ok := compareMap[player.PlayerId]
        if ok {
            player.FrameIdx = frame.FrameIdx
            compareMap[player.PlayerId][1] = player
        }
        if !ok {
            compareMap[player.PlayerId] = make([]types.PlayerPosition, 2)
            player.FrameIdx = frame.FrameIdx
            compareMap[player.PlayerId][1] = player
        }
    }

    c := 0
    for playerId, values := range compareMap {

        xDiff := math.Abs(values[0].Position.X - values[1].Position.X)
        yDiff := math.Abs(values[0].Position.Y - values[1].Position.Y)
        if xDiff > JumpThreshold || yDiff > JumpThreshold {
            w.Logger.Info("Jump detected", "player_id", playerId, "x_diff", xDiff, "y_diff", yDiff, "frame", frame.FrameIdx)
            checkPlayer[fmt.Sprintf("%s:%d:0", playerId, frame.FrameIdx)] = values[0]
            checkPlayer[fmt.Sprintf("%s:%d:1", playerId, frame.FrameIdx)] = values[1]
        }
        if c == len(compareMap)-1 && frame.FrameIdx%5 == 0 {
            if len(checkPlayer) > 1 {
                w.swap(checkPlayer)
            }
            checkPlayer = make(map[string]types.PlayerPosition)
        }
        c++
    }
}

func (w *Worker) detectHoles(frame types.SmallFrame) {
    prevFrame, err := w.StateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - 1))
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
            var playerNum string // Placeholder: replace with actual PlayerNumber logic

            msgData := holeMessage{
                FrameIdx:     frame.FrameIdx,
                PlayerNumber: playerNum, // Use the determined PlayerNumber here
                Home:         false,
            }

            message, err := json.Marshal(msgData)
            if err != nil {
                w.Logger.Error("Failed to marshal holeMessage to JSON", "error", err, "playerId", playerId)
                continue // Skip to the next player if marshalling fails
            }

            // Declare messageLength first, then assign to existing err
            var messageLength int
            messageLength, err = w.socketImputation.SendMessage("hole", message) // Assuming topic is "hole"
            if err != nil {
                // Use messageLength in the error log
                w.Logger.Error("Failed to send hole message", "length", messageLength, "error", err, "playerId", playerId)
            }
        }
    }

    // Close the existing socket if it exists
    if w.socketImputation != nil {
        err = w.socketImputation.Close()
        if err != nil {
            w.Logger.Error("Failed to close existing socket", "error", err)
            return
        }
    }
}
func (w *Worker) swap(p map[string]types.PlayerPosition) {
    w.Logger.Error("Swapping players", "players", p)
}
