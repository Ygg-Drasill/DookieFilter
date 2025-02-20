package main

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"time"
)

const (
	SCREEN_W = 1200
	SCREEN_H = 800
	SCALE    = 4

	FIELD_W = 105
	FIELD_H = 68
)

var (
	RED = color.RGBA{
		R: 255,
	}
	BLUE = color.RGBA{B: 255}
)

type Game struct {
	fr          *frameReader.FrameReader
	ball        []float64
	awayPlayers map[string]types.Player
	homePlayers map[string]types.Player
	time        time.Time
	done        bool
}

func (g *Game) Update() error {
	if g.done {
		return nil
	}
	frame := g.fr.Next()
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	debugInfo := fmt.Sprintf("%d:%d:%d", g.time.Hour(), g.time.Minute(), g.time.Second())

	xoff := float32(SCREEN_W / 2)
	yoff := float32(SCREEN_H / 2)

	x := xoff + float32(g.ball[0])*SCALE
	y := yoff + float32(g.ball[1])*SCALE

	vector.DrawFilledCircle(screen, x, y, 4, color.White, true)

	debugInfo = fmt.Sprintf("%s away:%d", debugInfo, len(g.awayPlayers))
	for _, p := range g.awayPlayers {
		px, py := float32(p.Xyz[0]), float32(p.Xyz[1])
		vector.DrawFilledCircle(screen, px*SCALE+xoff, py*SCALE+yoff, 8, RED, true)
	}

	debugInfo = fmt.Sprintf("%s home:%d", debugInfo, len(g.homePlayers))
	for _, p := range g.homePlayers {
		px, py := float32(p.Xyz[0]), float32(p.Xyz[1])
		vector.DrawFilledCircle(screen, px*SCALE+xoff, py*SCALE+yoff, 8, BLUE, true)
	}

	if g.done {
		debugInfo = fmt.Sprintf("%s %s", debugInfo, "done")
	}
	ebitenutil.DebugPrint(screen, debugInfo)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {

	fr := frameReader.New("./raw.jsonl")

	ebiten.SetWindowSize(SCREEN_W, SCREEN_H)
	ebiten.SetWindowTitle("Visunator")
	if err := ebiten.RunGame(&Game{fr: fr, done: false, awayPlayers: make(map[string]types.Player), homePlayers: make(map[string]types.Player)}); err != nil {
		log.Fatal(err)
	}
}
