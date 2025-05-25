package app

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"visunator/font"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.width, g.height = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.frameIndex == -1 {
		return
	}
	debugInfo := fmt.Sprintf("%d:%d:%d", g.time.Hour(), g.time.Minute(), g.time.Second())

	xoff := float32(g.width / 2)
	yoff := float32(g.height / 2)

	var (
		x float32 = 0
		y float32 = 0
	)
	if g.ball != nil {
		x = xoff + float32(g.ball[0])*g.scale
		y = yoff + float32(g.ball[1])*g.scale
	}

	vector.DrawFilledCircle(screen, x, y, BALL_SIZE, color.White, true)
	vector.DrawFilledRect(screen, 0, float32(g.height)-12,
		float32(g.frameIndex)/float32(g.frameLoader.FrameCount())*float32(g.width),
		float32(g.height), color.White, true)

	debugInfo = fmt.Sprintf("%s\naway:%d", debugInfo, len(g.awayPlayers))
	for _, p := range g.awayPlayers {
		px, py := float32(p.Xyz[0])*g.scale+xoff, float32(p.Xyz[1])*g.scale+yoff
		DrawPlayer(screen, p, RED, px, py)
	}

	debugInfo = fmt.Sprintf("%s\nhome:%d", debugInfo, len(g.homePlayers))
	for _, p := range g.homePlayers {
		px, py := float32(p.Xyz[0])*g.scale+xoff, float32(p.Xyz[1])*g.scale+yoff
		DrawPlayer(screen, p, BLUE, px, py)
	}

	debugInfo = fmt.Sprintf("%s\n%dms per frame", debugInfo, g.frameTime)
	if g.done {
		debugInfo = fmt.Sprintf("%s\n%s", debugInfo, "done")
	}
	ebitenutil.DebugPrint(screen, debugInfo)
	fieldX := xoff - FIELD_W*g.scale/2
	fieldY := yoff - FIELD_H*g.scale/2
	vector.StrokeRect(screen, fieldX, fieldY, FIELD_W*g.scale, FIELD_H*g.scale, 1, color.White, false)
}

func DrawPlayer(screen *ebiten.Image, player types.Player, col color.Color, x, y float32) {
	vector.DrawFilledCircle(screen, x, y, PLAYER_SIZE, col, true)
	text.Draw(screen, player.Number, font.VisunatorFontFace, int(x), int(y), color.White)
}
