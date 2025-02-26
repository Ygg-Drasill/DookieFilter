package font

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

var VisunatorFontFace font.Face

func init() {
	fontData := fonts.PressStart2P_ttf
	parsedFont, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}
	VisunatorFontFace, err = opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	log.Println("loaded font face", VisunatorFontFace)
}
