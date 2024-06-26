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
	"github.com/yngvark/ebit-grid-test/pkg/camera"
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

// Tile types.
const (
	Grass = iota
	Water
	Mountain
)

// TileSize is the size of a tile in pixels.
const TileSize = 32

type Game struct {
	grassImage      *ebiten.Image
	mountainImage   *ebiten.Image
	waterImage      *ebiten.Image
	debugBackground *ebiten.Image

	middleRectangle *ebiten.Image
	rectangleX      float64
	rectangleY      float64

	gopherImage *ebiten.Image

	worldMap *world_map.WorldMap

	audioPlayer *audio.Player

	soundTicker *time.Ticker

	// The viewport camera
	camera *camera.Camera
	world  *ebiten.Image

	screenWidth  int
	screenHeight int
	layoutInited bool

	mapViewport          *image.Rectangle
	mapStartCoordinate   Coordinate
	currentMapCoordinate Coordinate
}

const mapMoveSpeed = 40
const zoomSpeed = 3
const rotationSpeed = 5
const edgeThreshold = 50

func (g *Game) Update() error {
	// -------------------------------------------------------------------------------------
	// Map movement with mouse
	// -------------------------------------------------------------------------------------

	/*
		// Get the current mouse position
		mouseX, mouseY := ebiten.CursorPosition()

		// Check if the mouse is near the left edge of the screen
		if mouseX <= edgeThreshold {
			g.camera.Position[0] -= mapMoveSpeed
		}

		// Check if the mouse is near the right edge of the screen
		if mouseX >= g.screenWidth-edgeThreshold {
			g.camera.Position[0] += mapMoveSpeed
		}

		// Check if the mouse is near the top edge of the screen
		if mouseY <= edgeThreshold {
			g.camera.Position[1] -= mapMoveSpeed
		}

		// Check if the mouse is near the bottom edge of the screen
		if mouseY >= g.screenHeight-edgeThreshold {
			g.camera.Position[1] += mapMoveSpeed
		}
	*/
	// -------------------------------------------------------------------------------------
	// Map movement with keyboard
	// -------------------------------------------------------------------------------------
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.camera.Position[0] -= mapMoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.camera.Position[0] += mapMoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.camera.Position[1] -= mapMoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.camera.Position[1] += mapMoveSpeed
	}

	_, wy := ebiten.Wheel()
	if ebiten.IsKeyPressed(ebiten.KeyQ) || wy < 0 {
		if g.camera.ZoomFactor > -2400 {
			g.camera.ZoomFactor -= zoomSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) || wy > 0 {
		if g.camera.ZoomFactor < 2400 {
			g.camera.ZoomFactor += zoomSpeed
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.camera.Rotation += rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.camera.Rotation -= rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.camera.Reset()
		g.currentMapCoordinate = g.mapStartCoordinate
		g.mapViewport = getMapViewportFromCoordinate(g.screenWidth, g.screenHeight, g.currentMapCoordinate)
	}

	// -------------------------------------------------------------------------------------
	// Other stuff
	// -------------------------------------------------------------------------------------
	g.rectangleX += float64(rand.Intn(3)) - 1
	g.rectangleY += float64(rand.Intn(3)) - 1

	var xMovement int
	var yMovement int
	recaulcateMapViewport := false

	if g.camera.Position[0] >= TileSize {
		xMovement = 1
		recaulcateMapViewport = true
		g.camera.Position[0] = 0
	} else if g.camera.Position[0] <= -TileSize {
		xMovement = -1
		recaulcateMapViewport = true
		g.camera.Position[0] = 0
	}

	if g.camera.Position[1] >= TileSize {
		yMovement = 1
		recaulcateMapViewport = true
		g.camera.Position[1] = 0
	} else if g.camera.Position[1] <= -TileSize {
		yMovement = -1
		recaulcateMapViewport = true
		g.camera.Position[1] = 0
	}

	// TILTAK:
	// - Jeg må tegne minst 1 ekstra tile i hver retning, slik at hvitt ikke vises når man beveger seg. Ta enda flere hvis jeg
	// er redd for ytelse.
	// - Sjekk camera position vs world pos eller hva det hetter. Må sjekke om transformasjonene kan ødelegge for noe.

	// When camera position has moved beyound a tile, we want to re-calculate the map viewport (which part of the map
	// is shown on the screen). And reset the camera position.
	// In other words, we don't move the camera position around a huge map, we just allow moving the camera within
	// the tile size.
	if recaulcateMapViewport {
		g.currentMapCoordinate = NewCoordinate(
			g.currentMapCoordinate.X+xMovement,
			g.currentMapCoordinate.Y+yMovement)

		g.mapViewport = getMapViewportFromCoordinate(g.screenWidth, g.screenHeight, g.currentMapCoordinate)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.DrawTiles(g.world, g.worldMap)
	g.camera.Render(g.world, screen)
	g.drawMovingRectangle(screen)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f\nMove (WASD/Arrows)\nZoom (QE)\nRotate (ZC)\nReset (Space)", ebiten.ActualTPS()),
	)

	ebitenX, ebitenY := ebiten.CursorPosition()
	worldX, worldY := g.camera.ScreenToWorld(ebitenX, ebitenY)

	// Draw debug info
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.screenWidth/2-50), float64(g.screenHeight-270))
	screen.DrawImage(g.debugBackground, op)

	debugInfo :=
		fmt.Sprintln("Custom version: hei1") +
			fmt.Sprintln(g.camera.String()) +
			fmt.Sprintf("Cursor world pos: %.2f,%.2f\n", worldX, worldY) +
			fmt.Sprintf("Cursor ebiten pos: %d,%d\n", ebitenX, ebitenY) +
			fmt.Sprintf("g.camera.Position[0]: %.2f\n", g.camera.Position[0]) +
			fmt.Sprintf("g.camera.Position[1]: %.2f\n", g.camera.Position[1]) +
			fmt.Sprintf("currentMapCoordinate: %s\n", g.currentMapCoordinate.String()) +
			fmt.Sprintf("mapStartCoordinate: %s\n", g.mapStartCoordinate.String()) +
			fmt.Sprintf("mapViewport: %s\n", g.mapViewport.String())

	ebitenutil.DebugPrintAt(screen, debugInfo, g.screenWidth/2, g.screenHeight-250)
}

func (g *Game) DrawTiles(screen *ebiten.Image, worldMap *world_map.WorldMap) {
	// Get coordinates of the world map to draw. The viewport only shows a part of the world map, so we only draw that.
	// cameraViewport.Min.X is the x position of the viewport, similar for y.
	xCoordMin := g.mapViewport.Min.X
	xCoordMax := g.mapViewport.Max.X

	yCoordMin := g.mapViewport.Min.Y
	yCoordMax := g.mapViewport.Max.Y

	x := 0
	y := 0

	for yCoord := yCoordMin; yCoord < yCoordMax; yCoord++ {
		for xCoord := xCoordMin; xCoord <= xCoordMax; xCoord++ {
			tile := worldMap.TileAt(xCoord, yCoord)

			tileType := tile % 3
			var tileImage *ebiten.Image
			switch tileType {
			case Grass:
				tileImage = g.grassImage
			case Water:
				tileImage = g.waterImage
			case Mountain:
				tileImage = g.mountainImage
			default:
				panic("Unknown tile type")
			}

			// Draw tile at x, y
			op := &ebiten.DrawImageOptions{}
			// Set the image's pixel position
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(tileImage, op)

			x += TileSize
		}

		y += TileSize
		x = 0
	}
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
	if outsideHeight != g.screenHeight || outsideWidth != g.screenWidth {
		log.Print("Window resized to ", outsideWidth, "x", outsideHeight)
		g.screenWidth = outsideWidth
		g.screenHeight = outsideHeight

		if !g.layoutInited {
			g.camera = camera.NewCamera(outsideWidth, outsideHeight)

			g.mapStartCoordinate = NewCoordinate(g.worldMap.Width()/2, g.worldMap.Height()/2)
			g.currentMapCoordinate = g.mapStartCoordinate

			g.mapViewport = getMapViewportFromCoordinate(outsideWidth, outsideHeight, g.mapStartCoordinate)

			g.world = ebiten.NewImage(outsideWidth, outsideHeight)

			g.layoutInited = true
		}
	}

	//g.tilesDrawer.Layout(outsideWidth, outsideHeight, g.worldMap)

	return outsideWidth, outsideHeight
}

func NewGame() (*Game, error) {
	g := &Game{}

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

	// Tiles
	g.grassImage = ebiten.NewImage(TileSize, TileSize)
	g.grassImage.Fill(color.NRGBA{R: 188, G: 231, B: 132, A: 0xff})

	g.mountainImage = ebiten.NewImage(TileSize, TileSize)
	g.mountainImage.Fill(color.NRGBA{R: 66, G: 62, B: 55, A: 0xff})

	g.waterImage = ebiten.NewImage(TileSize, TileSize)
	g.waterImage.Fill(color.NRGBA{R: 52, G: 138, B: 167, A: 0xff})

	g.debugBackground = ebiten.NewImage(350, 250)
	g.debugBackground.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: 0xff})

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
	worldMap := world_map.Generate(1500, 1500, 600)
	g.worldMap = worldMap

	return g, nil
}
