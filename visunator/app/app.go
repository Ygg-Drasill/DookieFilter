package app

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"time"
	"visunator/font"
	_ "visunator/font"
)

const (
	SCREEN_W = 1200
	SCREEN_H = 800
	SCALE    = 8

	FIELD_W = 105
	FIELD_H = 68

	PLAYER_SIZE = 8
	BALL_SIZE   = 4
)

var (
	RED = color.RGBA{
		R: 255,
	}
	BLUE = color.RGBA{B: 255}
)

type Game struct {
	frameLoader types.FrameLoader[types.DataPlayer]
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
	frame := g.frameLoader.Next()
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

func DrawPlayer(screen *ebiten.Image, player types.Player, col color.Color, x, y float32) {
	vector.DrawFilledCircle(screen, x, y, PLAYER_SIZE, col, true)
	text.Draw(screen, player.Number, font.VisunatorFontFace, int(x), int(y), color.White)
}

func (g *Game) Draw(screen *ebiten.Image) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, _ := ebiten.CursorPosition()
		windowWidth := 1000
		frameIndex := int64(float32(mx) / float32(windowWidth) * float32(g.frameLoader.FrameCount()))
		fmt.Println(frameIndex)
		err := g.frameLoader.GoToFrame(frameIndex)
		if err != nil {
			log.Fatal(err)
		}
	}

	debugInfo := fmt.Sprintf("%d:%d:%d", g.time.Hour(), g.time.Minute(), g.time.Second())

	xoff := float32(SCREEN_W / 2)
	yoff := float32(SCREEN_H / 2)

	x := xoff + float32(g.ball[0])*SCALE
	y := yoff + float32(g.ball[1])*SCALE

	vector.DrawFilledCircle(screen, x, y, BALL_SIZE, color.White, true)

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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Run() {
	ebiten.SetWindowSize(SCREEN_W, SCREEN_H)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Visunator")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func NewFromReader(path string) *Game {
	fr := frameReader.New(path)
	return &Game{
		frameLoader: fr,
		done:        false,
		awayPlayers: make(map[string]types.Player),
		homePlayers: make(map[string]types.Player),
	}
}
