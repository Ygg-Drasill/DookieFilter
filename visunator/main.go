package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"time"
	"visunator/frameReader"
	"visunator/types"
)

const (
	SCREEN_W = 640
	SCREEN_H = 480
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

	x := SCREEN_W/2 + float32(g.ball[0])
	y := SCREEN_H/2 + float32(g.ball[1])

	vector.DrawFilledCircle(screen, x, y, 8, color.White, true)

	debugInfo = fmt.Sprintf("%s away:%d", debugInfo, len(g.awayPlayers))
	for _, p := range g.awayPlayers {
		px, py := float32(p.Xyz[0]), float32(p.Xyz[1])
		vector.DrawFilledCircle(screen, px, py, 16, RED, true)
	}

	debugInfo = fmt.Sprintf("%s home:%d", debugInfo, len(g.homePlayers))
	for _, p := range g.homePlayers {
		px, py := float32(p.Xyz[0]), float32(p.Xyz[1])
		vector.DrawFilledCircle(screen, px, py, 16, RED, true)
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
