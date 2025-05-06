package detector

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"math"
	"strconv"
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
			checkPlayer[fmt.Sprintf("%s:%d:0", playerId, prevFrame.FrameIdx)] = values[0]
			checkPlayer[fmt.Sprintf("%s:%d:1", playerId, frame.FrameIdx)] = values[1]
		}
		if c == len(compareMap)-1 && frame.FrameIdx%5 == 0 {
			if len(checkPlayer) > 1 {
				w.decide(checkPlayer)
			}
		}
		c++
	}

}

func (w *Worker) decide(p map[string]types.PlayerPosition) {
	x := len(p)
	switch x {
	case 2: // one player
		w.jump(p)
	default:
		parity := x%2 == 0
		w.swap(parity, p)
	}
}

type swapPlayer struct {
	key    string
	player types.PlayerPosition
}

func (w *Worker) swap(parity bool, p map[string]types.PlayerPosition) {
	var (
		cf,
		pf []swapPlayer
	)
	for k, player := range p {
		if player.FrameIdx == 0 {
			w.Logger.Debug("hole", "key", k, "player", player)
			q := strings.Split(k, ":")
			s, err := strconv.Atoi(q[1])
			if err != nil {
				w.Logger.Error("Failed to parse frame index", "error", err)
				continue
			}
			player.FrameIdx = s
			player.PlayerId = q[0]
			p[k] = player
		}
		if strings.HasSuffix(k, ":0") {
			pf = append(pf, swapPlayer{key: k, player: player})
		}
		if strings.HasSuffix(k, ":1") {
			cf = append(cf, swapPlayer{key: k, player: player})
		}
	}

	if len(pf) != len(cf) {
		w.Logger.Error("swap", "error", "Mismatched frame data", "pf", pf, "cf", cf)
		return
	}

	for _, prev := range pf {
		for _, curr := range cf {
			if prev.player.PlayerId != curr.player.PlayerId && positionProximity(prev, curr) {
				swapPlayers(p, curr, getPair(cf, prev))
			}
		}
	}
}

func getPair(players []swapPlayer, m swapPlayer) swapPlayer {
	for _, p := range players {
		if p.key != m.key && p.player.PlayerId == m.player.PlayerId {
			return p
		}
	}
	return m
}

const (
	fieldLength = 105.0
	fieldWidth  = 68.0

	playerMaxSpeed = 9.5
	frameTime      = 1.0 / 25.0
)

var (
	fieldSize      = math.Sqrt(fieldLength*fieldLength + fieldWidth*fieldWidth)
	maxMovePercent = frameTime * playerMaxSpeed
)

func positionProximity(p1, p2 swapPlayer) bool {
	dx := p1.player.Position.X - p2.player.Position.X
	dy := p1.player.Position.Y - p2.player.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	normalized := (distance / fieldSize) * 100
	fmt.Println("Normalized distance (%)", normalized)

	return normalized < maxMovePercent
}

func swapPlayers(
	p map[string]types.PlayerPosition,
	p1, p2 swapPlayer) {
	tmpI := p1.player
	tmpJ := p2.player
	p1.player.Position = tmpJ.Position
	p2.player.Position = tmpI.Position
	p[p1.key] = p1.player
	p[p2.key] = p2.player
}

func (w *Worker) jump(p map[string]types.PlayerPosition) bool {
	var pf, cf types.PlayerPosition

	for k, player := range p {
		if strings.HasSuffix(k, ":0") {
			pf = player
		}
		if strings.HasSuffix(k, ":1") {
			cf = player
		}
	}

	if pf.FrameIdx == 0 {
		w.Logger.Warn("hole", "key", pf.PlayerId, "player", pf)
		return false
	}
	if cf.FrameIdx == 0 {
		w.Logger.Warn("hole", "key", cf.PlayerId, "player", cf)
		return false
	}

	prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(pf.FrameIdx))
	if err != nil {
		w.Logger.Warn("No previous frame to compare", "frame", pf.FrameIdx)
		return false
	}

	for _, player := range prevFrame.Players {
		if player.PlayerId != cf.PlayerId &&
			positionsEqual(player.Position, cf.Position) {
			w.Logger.Info("Duplicate player coordinates detected",
				"jump_player", cf,
				"real_player", player)
			return true
		}
	}
	return false
}

func positionsEqual(p1, p2 types.Position) bool {
	const e = 0.01
	return math.Abs(p1.X-p2.X) < e &&
		math.Abs(p1.Y-p2.Y) < e
}
