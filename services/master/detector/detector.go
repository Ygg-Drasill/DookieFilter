package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
    "github.com/Ygg-Drasill/DookieFilter/common/types"
    "github.com/Ygg-Drasill/DookieFilter/services/master/worker"
    zmq "github.com/pebbe/zmq4"
    "strings"
    "sync"
)

type gameState struct {
    players  []types.PlayerPosition
    ball     types.PlayerPosition
    frameIdx int
}

func (g gameState) Key() pringleBuffer.Key {
    return pringleBuffer.Key(g.frameIdx)
}

type DetectorWorker struct {
    worker.BaseWorker

    socketListen *zmq.Socket

    stateBuffer *pringleBuffer.PringleBuffer[gameState]
}

func New(ctx *zmq.Context, options ...func(worker *DetectorWorker)) *DetectorWorker {
    w := &DetectorWorker{
        BaseWorker:  worker.NewBaseWorker(ctx, "detector"),
        stateBuffer: pringleBuffer.New[gameState](10),
    }
    for _, opt := range options {
        opt(w)
    }

    return w
}

func (w *DetectorWorker) Run(wg *sync.WaitGroup) {
    defer wg.Done()
    defer w.close()
    w.Logger.Info("Starting detector worker")
    err := w.connect()
    if err != nil {
        w.Logger.Error("Failed to bind/connect zmq sockets", "error", err)
    }

    for {
        topic, _ := w.socketListen.Recv(zmq.SNDMORE)
        if topic == "playerframe" {
            message, _ := w.socketListen.RecvMessage(0)
            players := types.DeserializePlayerPositions(strings.Join(message, ""))
            w.stateBuffer.Insert(gameState{players: players, frameIdx: players[0].FrameIdx})
        }
    }
}
