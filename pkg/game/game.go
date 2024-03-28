package game

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rs/zerolog/log"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles/world_map"
	"image"
	"image/color"
	"math/rand"
	"time"
)

//go:embed gopher.png
var gopherPng []byte

//go:embed Glass.mp3
var soundBytes []byte

type Game struct {
	tilesDrawer *tiles.Drawer

	gopherImage     *ebiten.Image
	middleRectangle *ebiten.Image

	rectangleX float64
	rectangleY float64

	worldMap [][]int

	outsideWidth  int
	outsideHeight int

	audioPlayer *audio.Player

	soundTicker *time.Ticker
}

func (g *Game) Update() error {
	g.rectangleX += float64(rand.Intn(3)) - 1
	g.rectangleY += float64(rand.Intn(3)) - 1

	select {
	case <-g.soundTicker.C:
		err := g.audioPlayer.Rewind()
		if err != nil {
			log.Err(err).Msg("rewinding audio player")
			break
		}

		g.audioPlayer.Play()

	default:
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	ebitenutil.DebugPrint(screen, "Hello, World!")

	g.drawStillImage(screen)
	g.tilesDrawer.Draw(screen, g.worldMap)
	g.drawMovingRectangle(screen)

	fmt.Printf("Screen size: %dx%d\n", screen.Bounds().Dx(), screen.Bounds().Dy())
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
	if outsideHeight != g.outsideHeight || outsideWidth != g.outsideWidth {
		log.Print("Window resized to ", outsideWidth, "x", outsideHeight)
		g.outsideWidth = outsideWidth
		g.outsideHeight = outsideHeight
	}

	return outsideWidth, outsideHeight
}

func NewGame() (*Game, error) {
	g := &Game{
		tilesDrawer: tiles.NewDrawer(),
	}

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Gopher
	gopherImageSource, _, err := image.Decode(bytes.NewReader(gopherPng))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	g.gopherImage = ebiten.NewImageFromImage(gopherImageSource)

	// Middle rectangle
	g.middleRectangle = ebiten.NewImage(50, 50)
	g.middleRectangle.Fill(color.NRGBA{R: 0x80, G: 0, B: 0, A: 0xff})

	// Sound
	audioContext := audio.NewContext(44100)

	soundStream, err := mp3.DecodeWithSampleRate(audioContext.SampleRate(), bytes.NewReader(soundBytes))
	if err != nil {
		return nil, fmt.Errorf("decoding audio: %w", err)
	}

	audioPlayer, err := audioContext.NewPlayer(soundStream)
	if err != nil {
		return nil, fmt.Errorf("creating audio player: %w", err)
	}

	g.audioPlayer = audioPlayer

	// Ticker
	g.soundTicker = time.NewTicker(2 * time.Second)

	// World map
	g.worldMap = world_map.Generate(20, 10, 3)
	return g, nil
}
