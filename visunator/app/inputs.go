package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

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

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.frameTime++
	} else if ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.frameTime++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.frameTime--
	} else if ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.frameTime--
	}

	if g.frameTime < 20 {
		g.frameTime = 20
	}
}
