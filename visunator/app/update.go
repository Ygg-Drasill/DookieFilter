package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"time"
)

func (g *Game) Update() error {
	g.HandleInputs()
	if !g.active || g.done && g.frameIndex != -1 {
		return nil
	}
	tNow := time.Now().UnixMilli()
	tSinceLastUpdate := tNow - g.lastUpdate
	if int(tSinceLastUpdate) < g.updateFrequency {
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

	for _, p := range data.AwayPlayers {
		g.awayPlayers[p.PlayerId] = p
	}

	for _, p := range data.HomePlayers {
		g.homePlayers[p.PlayerId] = p
	}

	g.lastUpdate = tNow
	return nil
}

func (g *Game) HandleInputs() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if my < g.height/2 {
			return
		}
		frameIndex := int64(float32(mx) / float32(g.width) * float32(g.frameLoader.FrameCount()))
		err := g.frameLoader.GoToFrame(frameIndex)
		g.frameIndex = frameIndex
		if err != nil {
			log.Fatal(err)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.active = !g.active
	}
}
