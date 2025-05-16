package detector

import (
	"errors"
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

const (
	fieldLength = 105.0
	fieldWidth  = 68.0

	playerMaxSpeed = 9.5 // m/s
	frameTime      = 1.0 / 25.0
)

var (
	fieldSize       = math.Sqrt(fieldLength*fieldLength + fieldWidth*fieldWidth)
	maxMovePerFrame = frameTime * playerMaxSpeed
)

func (w *Worker) detect(frame types.SmallFrame) {
	prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - 1))
	if errors.Is(err, pringleBuffer.PringleBufferError{}) {
		w.Logger.Warn("No previous frame to compare")
		return
	}
	if err != nil {
		w.Logger.Error("Failed to get previous frame", "error", err)
		return
	}

	compareMap := make(map[int][]types.PlayerPosition)
	for _, player := range prevFrame.Players {
		_, ok := compareMap[player.PlayerNum]
		if ok {
			player.FrameIdx = prevFrame.FrameIdx
			compareMap[player.PlayerNum][0] = player
		}
		if !ok {
			compareMap[player.PlayerNum] = make([]types.PlayerPosition, 2)
			player.FrameIdx = prevFrame.FrameIdx
			compareMap[player.PlayerNum][0] = player
		}
	}
	for _, player := range frame.Players {
		_, ok := compareMap[player.PlayerNum]
		if ok {
			player.FrameIdx = frame.FrameIdx
			compareMap[player.PlayerNum][1] = player
		}
		if !ok {
			compareMap[player.PlayerNum] = make([]types.PlayerPosition, 2)
			player.FrameIdx = frame.FrameIdx
			compareMap[player.PlayerNum][1] = player
		}
	}

	var (
		count       = 0
		checkPlayer = make(map[string]types.PlayerPosition)
	)
	for playerId, values := range compareMap {
		xDiff := math.Abs(values[0].Position.X - values[1].Position.X)
		yDiff := math.Abs(values[0].Position.Y - values[1].Position.Y)
		if xDiff > maxMovePerFrame || yDiff > maxMovePerFrame {
			w.Logger.Info("Jump detected", "player_id", playerId, "x_diff", xDiff, "y_diff", yDiff, "frame", frame.FrameIdx)
			checkPlayer[fmt.Sprintf("%s:%d:0", playerId, prevFrame.FrameIdx)] = values[0]
			checkPlayer[fmt.Sprintf("%s:%d:1", playerId, frame.FrameIdx)] = values[1]
		}
		if count == len(compareMap)-1 && frame.FrameIdx%5 == 0 {
			if len(checkPlayer) > 1 {
				w.decide(w.swap(checkPlayer), checkPlayer, &frame)
				w.stateBuffer.Insert(frame)
			}
		}
		count++
	}
}

func (w *Worker) decide(
	swappers map[string]bool,
	p map[string]types.PlayerPosition,
	frame *types.SmallFrame) *types.SmallFrame {

	for key, swapped := range swappers {
		playerId := strings.Split(key, ":")[0]
		found := false
		for i, f := range frame.Players {
			if f.PlayerId != playerId {
				continue
			}
			found = true
			if !swapped {
				q, err := w.stateBuffer.Get(pringleBuffer.Key(f.FrameIdx - 1))
				if err != nil {
					w.Logger.Error("Failed to get previous frame", "error", err)
					continue
				}
				if clearPlayer(playerId, q) {
					break
				}
				frame.Players[i].Position = types.Position{X: 0, Y: 0}
				break
			}
			frame.Players[i].Position = p[key].Position
			swappers[key] = false
			w.Logger.Debug("swapped", "key", key, "player", f)
			break
		}
		if swapped && !found {
			addPlayer(frame, p[key])
		}
	}

	return frame
}

func clearPlayer(playerId string, q types.SmallFrame) bool {
	for _, player := range q.Players {
		if player.PlayerId == playerId &&
			player.Position.X == 0 && player.Position.Y == 0 {
			return true
		}
	}
	return false
}

func addPlayer(frame *types.SmallFrame, player types.PlayerPosition) {
	frame.Players = append(frame.Players, player)
}

type swapPlayer struct {
	key    string
	player types.PlayerPosition
}

func (w *Worker) swap(p map[string]types.PlayerPosition) (swapped map[string]bool) {
	var (
		cf,
		pf []swapPlayer
		swappers = make(map[string]bool)
	)

	for k, player := range p {
		if player.FrameIdx == 0 {
			w.Logger.Debug("hole", "key", k, "player", player)
			q := strings.Split(k, ":")
			f, err := strconv.Atoi(q[1])
			if err != nil {
				w.Logger.Error("Failed to parse frame index", "error", err)
				continue
			}
			player.FrameIdx = f
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
		return swappers
	}

	// load swappers
	for _, s := range cf {
		swappers[s.key] = false
	}
	for _, prev := range pf {
		for _, curr := range cf {
			if prev.player.PlayerId != curr.player.PlayerId && positionProximity(prev, curr) {
				g := getPair(cf, prev)
				if swappers[curr.key] == true || swappers[g.key] == true {
					continue
				}
				swapPlayers(p, curr, g)
				swappers[curr.key] = true
				swappers[g.key] = true
				break
			}
		}
	}

	return swappers
}

func getPair(players []swapPlayer, m swapPlayer) swapPlayer {
	for _, p := range players {
		if p.key != m.key && p.player.PlayerId == m.player.PlayerId {
			return p
		}
	}
	return m
}

func positionProximity(p1, p2 swapPlayer) bool {
	dx := p1.player.Position.X - p2.player.Position.X
	dy := p1.player.Position.Y - p2.player.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	normalized := (distance / fieldSize) * 100

	return normalized < maxMovePerFrame
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
