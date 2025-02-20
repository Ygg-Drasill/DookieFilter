package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, 32, 32, 16, color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	data, err := os.ReadFile("../data/raw-data.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	fmt.Println(lines)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Visunator")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

type RawFrame struct {
	League    string `json:"league"`
	GameId    string `json:"gameId"`
	FeedName  string `json:"feedName"`
	MessageId string `json:"messageId"`
	Data      []struct {
		AwayPlayers []struct {
			Number   string    `json:"number"`
			OptaId   string    `json:"optaId"`
			PlayerId string    `json:"playerId"`
			Speed    float64   `json:"speed"`
			Xyz      []float64 `json:"xyz"`
		} `json:"awayPlayers"`
		Ball struct {
			Speed float64   `json:"speed"`
			Xyz   []float64 `json:"xyz"`
		} `json:"ball"`
		FrameIdx    int     `json:"frameIdx"`
		GameClock   float64 `json:"gameClock"`
		HomePlayers []struct {
			Number   string    `json:"number"`
			OptaId   string    `json:"optaId"`
			PlayerId string    `json:"playerId"`
			Speed    float64   `json:"speed"`
			Xyz      []float64 `json:"xyz"`
		} `json:"homePlayers"`
		Period    int   `json:"period"`
		WallClock int64 `json:"wallClock"`
	} `json:"data"`
}

type ProducedFrame struct {
	Period      int     `json:"period"`
	FrameIdx    int     `json:"frameIdx"`
	GameClock   float64 `json:"gameClock"`
	WallClock   int64   `json:"wallClock"`
	HomePlayers []struct {
		PlayerId string    `json:"playerId"`
		Number   int       `json:"number"`
		Xyz      []float64 `json:"xyz"`
		Speed    float64   `json:"speed"`
		OptaId   string    `json:"optaId"`
	} `json:"homePlayers"`
	AwayPlayers []struct {
		PlayerId string    `json:"playerId"`
		Number   int       `json:"number"`
		Xyz      []float64 `json:"xyz"`
		Speed    float64   `json:"speed"`
		OptaId   string    `json:"optaId"`
	} `json:"awayPlayers"`
	Ball struct {
		Xyz   []float64 `json:"xyz"`
		Speed float64   `json:"speed"`
	} `json:"ball"`
	Live      bool   `json:"live"`
	LastTouch string `json:"lastTouch"`
}
