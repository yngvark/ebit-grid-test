package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed gopher.png
var gopherPng []byte

// Images
var gopherImage *ebiten.Image
var middleRectangle *ebiten.Image

type Game struct {
	inited bool

	rectangleX float64
	rectangleY float64
}

func (g *Game) init() error {
	defer func() {
		g.inited = true
	}()

	// Gopher
	gopherImageSource, _, err := image.Decode(bytes.NewReader(gopherPng))
	if err != nil {
		return fmt.Errorf("decoding image: %w", err)
	}

	gopherImage = ebiten.NewImageFromImage(gopherImageSource)

	// Middle rectangle
	middleRectangle = ebiten.NewImage(50, 50)
	middleRectangle.Fill(color.NRGBA{R: 0x80, G: 0, B: 0, A: 0xff})

	tiles.Init()

	return nil
}

func (g *Game) Update() error {
	if !g.inited {
		err := g.init()
		if err != nil {
			return fmt.Errorf("initializing: %w", err)
		}
	}

	g.rectangleX += float64(rand.Intn(3)) - 1
	g.rectangleY += float64(rand.Intn(3)) - 1

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	ebitenutil.DebugPrint(screen, "Hello, World!")

	g.drawStillImage(screen)
	tiles.Draw(screen)
	g.drawMovingRectangle(screen)
}

func (g *Game) drawStillImage(screen *ebiten.Image) {
	screen.DrawImage(gopherImage, nil)
}

func (g *Game) drawMovingRectangle(screen *ebiten.Image) {
	//screenWidth, screenHeight := screen.Size()
	//rectangleWidth, rectangleHeight := middleRectangle.Size()

	// Calculate the x and y coordinates to draw the image at the center of the window.
	//x := float64(screenWidth/2-rectangleWidth/2) + g.rectangleX
	//y := float64(screenHeight/2-rectangleHeight/2) + g.rectangleY
	x := 50 + g.rectangleX
	y := 50 + g.rectangleY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	screen.DrawImage(middleRectangle, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	game := Game{}
	game.init()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
