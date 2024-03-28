package game

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles"
	"image"
	"image/color"
	"math/rand"
)

//go:embed gopher.png
var gopherPng []byte

type Game struct {
	tilesDrawer tiles.TilesDrawer

	gopherImage     *ebiten.Image
	middleRectangle *ebiten.Image

	rectangleX float64
	rectangleY float64
}

func (g *Game) init() error {
	// Gopher
	gopherImageSource, _, err := image.Decode(bytes.NewReader(gopherPng))
	if err != nil {
		return fmt.Errorf("decoding image: %w", err)
	}

	g.gopherImage = ebiten.NewImageFromImage(gopherImageSource)

	// Middle rectangle
	g.middleRectangle = ebiten.NewImage(50, 50)
	g.middleRectangle.Fill(color.NRGBA{R: 0x80, G: 0, B: 0, A: 0xff})

	return nil
}

func (g *Game) Update() error {
	g.rectangleX += float64(rand.Intn(3)) - 1
	g.rectangleY += float64(rand.Intn(3)) - 1

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	ebitenutil.DebugPrint(screen, "Hello, World!")

	g.drawStillImage(screen)
	g.tilesDrawer.Draw(screen)
	g.drawMovingRectangle(screen)
}

func (g *Game) drawStillImage(screen *ebiten.Image) {
	screen.DrawImage(g.gopherImage, nil)
}

func (g *Game) drawMovingRectangle(screen *ebiten.Image) {
	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	rectangleWidth := g.middleRectangle.Bounds().Dx()
	rectangleHeight := g.middleRectangle.Bounds().Dy()

	// Calculate the x and y coordinates to draw the image at the center of the window.
	x := float64(screenWidth/2-rectangleWidth/2) + g.rectangleX
	y := float64(screenHeight/2-rectangleHeight/2) + g.rectangleY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	screen.DrawImage(g.middleRectangle, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func NewGame() (*Game, error) {
	g := &Game{
		tilesDrawer: tiles.NewTiles(),
	}

	err := g.init()
	if err != nil {
		return nil, fmt.Errorf("initing: %w", err)
	}

	return g, nil
}
