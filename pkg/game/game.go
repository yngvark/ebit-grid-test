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
	"math"
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

	worldMap *world_map.WorldMap

	outsideWidth  int
	outsideHeight int

	audioPlayer *audio.Player

	soundTicker *time.Ticker

	// cameraViewport is the rectangle of the world map that is visible in the screen/window.
	cameraViewport *image.Rectangle
	scaleFactor    float64

	viewportInited bool
}

func (g *Game) Update() error {
	g.rectangleX += float64(rand.Intn(3)) - 1
	g.rectangleY += float64(rand.Intn(3)) - 1

	_, wy := ebiten.Wheel()
	if wy < 0 {
		g.scaleFactor *= 0.99 // Decrease the scale factor when scrolling down
	} else if wy > 0 {
		g.scaleFactor *= 1.01 // Increase the scale factor when scrolling up
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	ebitenutil.DebugPrint(screen, "Hello, World!")

	g.tilesDrawer.Draw(screen, g.worldMap, g.cameraViewport, g.scaleFactor)
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
	if outsideHeight != g.outsideHeight || outsideWidth != g.outsideWidth {
		log.Print("Window resized to ", outsideWidth, "x", outsideHeight)
		g.outsideWidth = outsideWidth
		g.outsideHeight = outsideHeight

		if !g.viewportInited {
			g.cameraViewport = getCameraViewportOfMapCenter(outsideWidth, outsideHeight, g.worldMap)
			g.viewportInited = true
		}
	}

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
	worldMap := world_map.Generate(700, 500, 600)
	g.worldMap = worldMap

	// Other
	g.scaleFactor = 1.0

	return g, nil
}

// getCameraViewportOfMapCenter calculates which part of the map should be visible inside the viewport.
func getCameraViewportOfMapCenter(windowWidth int, windowHeight int, worldMap *world_map.WorldMap) *image.Rectangle {
	windowWidthInCoords := int(math.Ceil(float64(windowWidth) / tiles.TileSize))
	windowHeightInCoords := int(math.Ceil(float64(windowHeight) / tiles.TileSize))

	// Use map center
	xCoordsMin := int(math.Floor(float64(worldMap.Width())/2 - float64(windowWidthInCoords)/2))
	yCoordsMin := int(math.Floor(float64(worldMap.Height()/2 - windowHeightInCoords/2)))

	xCoordsMax := int(math.Ceil(float64(xCoordsMin + windowWidthInCoords)))
	yCoordsMax := int(math.Ceil(float64(yCoordsMin + windowHeightInCoords)))

	viewPortCoords := image.Rect(
		xCoordsMin,
		yCoordsMin,
		xCoordsMax,
		yCoordsMax,
	)

	r := viewPortCoords
	return &r
}
