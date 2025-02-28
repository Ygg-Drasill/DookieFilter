package app

import (
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"time"
)

func (g *Game) Update() error {
	g.HandleInputs()
	if !g.active || g.done && g.frameIndex != -1 {
		return nil
	}
	tNow := time.Now().UnixMilli()
	tSinceLastUpdate := tNow - g.lastUpdate
	if int(tSinceLastUpdate) < g.frameTime {
		return nil
	}
	frame, err := g.frameLoader.Next()
	if err != nil {
		return err
	}
	g.frameIndex++
	if frame == nil {
		g.done = true
		return nil
	}
	data := frame.Data[0]
	g.ball = data.Ball.Xyz
	g.time = time.Unix(0, data.WallClock*int64(time.Millisecond))

	g.awayPlayers = make(map[string]types.Player)
	g.homePlayers = make(map[string]types.Player)

	for _, p := range data.AwayPlayers {
		g.awayPlayers[p.PlayerId] = p
	}

	for _, p := range data.HomePlayers {
		g.homePlayers[p.PlayerId] = p
	}

	g.lastUpdate = tNow
	return nil
}
