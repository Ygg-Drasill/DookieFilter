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
