package detector

import (
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

    socketListen  *zmq.Socket
    socketStorage *zmq.Socket

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
        }
    }
}

const JumpThreshold = 5 //TODO: change me

func (w *Worker) detect(frame types.SmallFrame) {
    var checkPlayer = make(map[string]types.PlayerPosition)
    prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - 1))
    if err != nil {
        w.Logger.Warn("No previous frame to compare")
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
        }
        c++
    }

}

const (
    onePlayer  = 2
    twoPlayers = 4
)

var evenPlayers = func(x int) bool { return x%2 == 0 }

func (w *Worker) swap(p map[string]types.PlayerPosition) {
    x := len(p)
    switch x {
    case onePlayer:
        fmt.Println("One player")
        w.jump(p)
    case twoPlayers:
        fmt.Println("Two players")
    default:
        if evenPlayers(x) {
            fmt.Printf("Even number of players: %d\n", x)
        } else {
            fmt.Printf("Odd number of players: %d\n", x)
        }
        fmt.Printf("Unknown number of players: %d\n", x)
    }
}

func (w *Worker) jump(p map[string]types.PlayerPosition) {
    var p1, p2 string
    for k, player := range p {
        if player.FrameIdx == 0 {
            fmt.Println("hole")
            return
        }
        if strings.HasSuffix(k, ":0") {
            p1 = k
        }
        if strings.HasSuffix(k, ":1") {
            p2 = k
        }
    }
    k := p[p1].FrameIdx
    s, err := w.stateBuffer.Get(pringleBuffer.Key(k))
    if err != nil {
        w.Logger.Warn("No previous frame to compare")
        return
    }
    for _, player := range s.Players {
        if player.Position == p[p1].Position && player.PlayerId != p[p1].PlayerId {
            fmt.Println("Jump detected", player.PlayerId)
            return
        }
        if player.Position == p[p2].Position && player.PlayerId != p[p2].PlayerId {
            fmt.Println("Jump detected", player.PlayerId)
            return
        }
        if player.Position == p[p1].Position {
            fmt.Println("Jump detected", p[p1].PlayerId)
            return
        }

    }
}
