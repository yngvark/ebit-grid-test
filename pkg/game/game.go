package game

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles/world_map"
	"image"
	"image/color"
	"time"
)

//go:embed gopher.png
var gopherPng []byte

//go:embed Glass.mp3
var soundBytes []byte

type Game struct {
	tilesDrawer *tiles.Drawer

	gopherImage *ebiten.Image

	worldMap *world_map.WorldMap

	audioPlayer *audio.Player

	soundTicker *time.Ticker
}

func (g *Game) Update() error {
	g.tilesDrawer.MoveRectangle()

	_, wy := ebiten.Wheel()
	if wy < 0 {
		g.tilesDrawer.DecreaseScaleFactor()
	} else if wy > 0 {
		g.tilesDrawer.IncreaseScaleFactor()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	ebitenutil.DebugPrint(screen, "Hello, World!")

	g.tilesDrawer.Draw(screen, g.worldMap)
	//g.drawMovingRectangle(screen)
}

func (g *Game) drawStillImage(screen *ebiten.Image) {
	screen.DrawImage(g.gopherImage, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.tilesDrawer.Layout(outsideWidth, outsideHeight, g.worldMap)

	return outsideWidth, outsideHeight
}

func NewGame() (*Game, error) {
	g := &Game{
		tilesDrawer: tiles.NewDrawer(),
	}

	//ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(false)
	ebiten.MaximizeWindow()

	// Gopher
	gopherImageSource, _, err := image.Decode(bytes.NewReader(gopherPng))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	g.gopherImage = ebiten.NewImageFromImage(gopherImageSource)

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
	worldMap := world_map.Generate(700, 500, 600)
	g.worldMap = worldMap

	return g, nil
}
