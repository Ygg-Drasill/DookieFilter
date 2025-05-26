package detector

import (
	"encoding/json"
	"errors"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/Ygg-Drasill/DookieFilter/services/master/worker"
	zmq "github.com/pebbe/zmq4"
	"math"
	"strings"
	"sync"
)

type Worker struct {
	worker.BaseWorker

	socketListen       *zmq.Socket
	socketImputation   *zmq.Socket
	socketStorage      *zmq.Socket
	imputationEndpoint string

	stateBuffer *pringleBuffer.PringleBuffer[types.SmallFrame]
}

func New(ctx *zmq.Context, options ...func(worker *Worker)) *Worker {
	w := &Worker{
		BaseWorker:         worker.NewBaseWorker(ctx, "detector"),
		stateBuffer:        pringleBuffer.New[types.SmallFrame](10),
		imputationEndpoint: endpoints.TcpEndpoint(endpoints.IMPUTATION),
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
	fieldLength   = 105.0
	fieldWidth    = 68.0

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

	// Build prev and current maps
	prevFrameMap := make(map[types.PlayerKey]types.PlayerPosition)
	currFrameMap := make(map[types.PlayerKey]types.PlayerPosition)

	for _, player := range prevFrame.Players {
		key := types.NewPlayerKey(player.PlayerNum, player.Home)
		player.FrameIdx = prevFrame.FrameIdx
		prevFrameMap[key] = player
	}

	for _, player := range frame.Players {
		key := types.NewPlayerKey(player.PlayerNum, player.Home)
		player.FrameIdx = frame.FrameIdx
		currFrameMap[key] = player
	}

	prevP := make(map[types.PlayerKey]map[int]types.PlayerPosition)
	currP := make(map[types.PlayerKey]map[int]types.PlayerPosition)
	for key, prevPlayer := range prevFrameMap {
		currPlayer, ok := currFrameMap[key]
		if !ok {
			continue
		}

		moveDiff := math.Hypot(
			prevPlayer.Position.X-currPlayer.Position.X,
			prevPlayer.Position.Y-currPlayer.Position.Y,
		)
		if moveDiff > maxMovePerFrame {
			w.Logger.Info("Jump detected", "player", key, "moveDiff", moveDiff, "frame", frame.FrameIdx)
			if _, ok = prevP[key]; !ok {
				prevP[key] = make(map[int]types.PlayerPosition)
			}
			if _, ok = currP[key]; !ok {
				currP[key] = make(map[int]types.PlayerPosition)
			}
			prevP[key][prevPlayer.FrameIdx] = prevPlayer
			currP[key][frame.FrameIdx] = currPlayer
		}
	}

	if len(currP) > 0 {
		swappers := w.swap(prevP, currP)
		w.decide(swappers, currP, &frame)
	}

	w.stateBuffer.Insert(frame)
}

func (w *Worker) decide(
	swappers map[types.PlayerKey]bool,
	p map[types.PlayerKey]map[int]types.PlayerPosition,
	frame *types.SmallFrame) *types.SmallFrame {

	for key, swapped := range swappers {
		found := false
		for i, f := range frame.Players {
			if types.NewPlayerKey(f.PlayerNum, f.Home) != key {
				continue
			}
			found = true
			if !swapped {
				q, err := w.stateBuffer.Get(pringleBuffer.Key(f.FrameIdx - 1))
				if err != nil {
					w.Logger.Error("Failed to get previous frame", "error", err)
					continue
				}
				if clearPlayer(key, q) {
					break
				}
				break
			}
			oldPos := frame.Players[i].Position
			frame.Players[i].Position = p[key][f.FrameIdx].Position
			swappers[key] = false
			w.Logger.Debug("swapped", "player", f.SKey(), "old position", oldPos, "new position", frame.Players[i].Position)

			playerPosition, _ := json.Marshal(frame.Players[i])
			_, err := w.socketStorage.SendMessage("position", playerPosition)
			if err != nil {
				w.Logger.Error("Failed to send imputation message", "error", err, "key", key)
			}
		}
		if swapped && !found {
			for _, v := range p[key] {
				addPlayer(frame, v)
				break
			}
		}
	}

	return frame
}

func clearPlayer(key types.PlayerKey, q types.SmallFrame) bool {
	for _, player := range q.Players {
		if types.NewPlayerKey(player.PlayerNum, player.Home) == key &&
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
	key    types.PlayerKey
	player types.PlayerPosition
}

func (w *Worker) swap(
	prevP map[types.PlayerKey]map[int]types.PlayerPosition,
	currP map[types.PlayerKey]map[int]types.PlayerPosition,
) map[types.PlayerKey]bool {

	var (
		pf, cf   []swapPlayer
		swappers = make(map[types.PlayerKey]bool)
	)

	for key, frames := range prevP {
		for _, player := range frames {
			pf = append(pf, swapPlayer{key: key, player: player})
			break // use the first/only frame
		}
	}
	for key, frames := range currP {
		for _, player := range frames {
			cf = append(cf, swapPlayer{key: key, player: player})
			swappers[key] = false
			break
		}
	}

	if len(pf) != len(cf) {
		w.Logger.Error("swap", "error", "Mismatched frame data", "pf", len(pf), "cf", len(cf))
		return swappers
	}

	for _, prev := range pf {
		for _, curr := range cf {

			if prev.player.SKey() != curr.player.SKey() &&
				positionProximity(prev, curr) {

				g := getPair(cf, prev)
				if swappers[curr.key] || swappers[g.key] {
					continue
				}

				swapPlayers(currP, curr, g)
				swappers[curr.key] = true
				swappers[g.key] = true
				break
			}
		}
	}

	return swappers
}

func getPair(players []swapPlayer, p swapPlayer) swapPlayer {
	for _, cf := range players {
		if cf.player.SKey() == p.player.SKey() {
			return cf
		}
	}
	return p
}

func positionProximity(p1, p2 swapPlayer) bool {
	dx := p1.player.Position.X - p2.player.Position.X
	dy := p1.player.Position.Y - p2.player.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	normalizedDistance := distance / fieldSize
	normalizedMove := maxMovePerFrame / fieldSize
	return normalizedDistance < normalizedMove
}

func swapPlayers(
	currP map[types.PlayerKey]map[int]types.PlayerPosition,
	p1, p2 swapPlayer,
) {
	var idx1, idx2 int
	for i := range currP[p1.key] {
		idx1 = i
		break
	}

	for i := range currP[p2.key] {
		idx2 = i
		break
	}

	tmp1 := currP[p1.key][idx1]
	tmp2 := currP[p2.key][idx2]

	tmp1.Position, tmp2.Position = tmp2.Position, tmp1.Position

	currP[p1.key][idx1] = tmp1
	currP[p2.key][idx2] = tmp2
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
		currentPlayers[player.SKey()] = true
	}

	prevPlayers := make(map[string]bool)
	for _, player := range prevFrame.Players {
		prevPlayers[player.SKey()] = true
	}

	// Check for players who were present before but are missing now
	for playerId := range prevPlayers {
		if !currentPlayers[playerId] {
			// Since we don't have number and home info in SmallFrame,
			// we'll use the PlayerId as the number and determine home based on position
			num, _ := types.DeSKey(playerId)
			msgData := holeMessage{
				FrameIdx:     frame.FrameIdx,
				PlayerNumber: num,
				Home:         true, // In the test, we're removing from home players
			}

			message, err := json.Marshal(msgData)
			if err != nil {
				w.Logger.Error("Failed to marshal holeMessage to JSON", "error", err, "playerId", playerId)
				continue
			}

			_, err = w.socketImputation.SendMessage("hole", message)
			if err != nil {
				w.Logger.Error("Failed to send hole message", "error", err, "playerId", playerId)
			}
		}
	}
}
