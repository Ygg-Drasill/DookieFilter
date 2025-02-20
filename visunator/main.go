package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"time"
	"visunator/frame"
	"visunator/frameReader"
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
	awayPlayers map[string]frame.Player
	homePlayers map[string]frame.Player
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
	x := SCREEN_W/2 + float32(g.ball[0])
	y := SCREEN_H/2 + float32(g.ball[1])

	vector.DrawFilledCircle(screen, x, y, 8, color.White, true)

	for _, p := range g.awayPlayers {
		px, py := float32(p.Xyz[0]), float32(p.Xyz[1])
		vector.DrawFilledCircle(screen, px, py, 16, RED, true)
	}

	if g.done {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%s - %s", g.time.String(), "done"))
	} else {
		ebitenutil.DebugPrint(screen, g.time.String())
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	fr := frameReader.New("./raw.jsonl")
	ebiten.SetWindowSize(SCREEN_W, SCREEN_H)
	ebiten.SetWindowTitle("Visunator")
	if err := ebiten.RunGame(&Game{fr: fr, done: false, awayPlayers: make(map[string]frame.Player), homePlayers: make(map[string]frame.Player)}); err != nil {
		log.Fatal(err)
	}
}
