package app

import (
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"math"
	"time"
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
	frameLoader     types.FrameLoader[types.DataPlayer]
	frameIndex      int64
	ball            []float64
	awayPlayers     map[string]types.Player
	homePlayers     map[string]types.Player
	time            time.Time
	done            bool
	active          bool
	width, height   int
	updateFrequency int
	lastUpdate      int64
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
		frameLoader:     fr,
		frameIndex:      -1,
		done:            false,
		active:          true,
		awayPlayers:     make(map[string]types.Player),
		homePlayers:     make(map[string]types.Player),
		updateFrequency: int(math.Floor(1000 / 25)),
		lastUpdate:      time.Now().UnixMilli(),
	}
}
