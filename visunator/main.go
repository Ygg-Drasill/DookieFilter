package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"visunator/frameReader"
)

const (
	SCREEN_W = 640
	SCREEN_H = 480
)

type Game struct {
	fr   *frameReader.FrameReader
	ball []float64
}

func (g *Game) Update() error {
	g.ball = g.fr.Next().Data[0].Ball.Xyz
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x := SCREEN_W/2 + float32(g.ball[0])
	y := SCREEN_H/2 + float32(g.ball[1])
	vector.DrawFilledCircle(screen, x, y, 16, color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	fr := frameReader.New("./raw.jsonl")
	ebiten.SetWindowSize(SCREEN_W, SCREEN_H)
	ebiten.SetWindowTitle("Visunator")
	if err := ebiten.RunGame(&Game{fr: fr}); err != nil {
		log.Fatal(err)
	}
}
