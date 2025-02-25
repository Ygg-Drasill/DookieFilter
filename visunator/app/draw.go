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

	xoff := float32(SCREEN_W / 2)
	yoff := float32(SCREEN_H / 2)

	x := xoff + float32(g.ball[0])*SCALE
	y := yoff + float32(g.ball[1])*SCALE

	vector.DrawFilledCircle(screen, x, y, BALL_SIZE, color.White, true)
	vector.DrawFilledRect(screen, 0, float32(g.height)-12,
		float32(g.frameIndex)/float32(g.frameLoader.FrameCount())*float32(g.width),
		float32(g.height), color.White, true)

	debugInfo = fmt.Sprintf("%s away:%d", debugInfo, len(g.awayPlayers))
	for _, p := range g.awayPlayers {
		px, py := float32(p.Xyz[0])*SCALE+xoff, float32(p.Xyz[1])*SCALE+yoff
		DrawPlayer(screen, p, RED, px, py)
	}

	debugInfo = fmt.Sprintf("%s home:%d", debugInfo, len(g.homePlayers))
	for _, p := range g.homePlayers {
		px, py := float32(p.Xyz[0])*SCALE+xoff, float32(p.Xyz[1])*SCALE+yoff
		DrawPlayer(screen, p, BLUE, px, py)
	}

	if g.done {
		debugInfo = fmt.Sprintf("%s %s", debugInfo, "done")
	}
	ebitenutil.DebugPrint(screen, debugInfo)
	fieldX := xoff - FIELD_W*SCALE/2
	fieldY := yoff - FIELD_H*SCALE/2
	vector.StrokeRect(screen, fieldX, fieldY, FIELD_W*SCALE, FIELD_H*SCALE, 1, color.White, false)
}

func DrawPlayer(screen *ebiten.Image, player types.Player, col color.Color, x, y float32) {
	vector.DrawFilledCircle(screen, x, y, PLAYER_SIZE, col, true)
	text.Draw(screen, player.Number, font.VisunatorFontFace, int(x), int(y), color.White)
}
