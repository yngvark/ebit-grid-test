package tiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles/world_map"
	"image"
	"image/color"
	"math/rand"
)

// Tile types.
const (
	Grass = iota
	Water
	Mountain
)

// TileSize is the size of a tile in pixels.
const TileSize = 32

type Drawer struct {
	grassImage    *ebiten.Image
	mountainImage *ebiten.Image
	waterImage    *ebiten.Image

	middleRectangle *ebiten.Image
	rectangleX      float64
	rectangleY      float64

	outsideWidth   int
	outsideHeight  int
	viewportInited bool
	mapViewport    *image.Rectangle
	screen         *ebiten.Image

	scaleFactor float64
	geomX       float64
}

func (d *Drawer) Draw(screen *ebiten.Image, worldMap *world_map.WorldMap) {
	// 1 tegn alle mine tiles på min eget image, world
	d.DrawTiles(d.screen, worldMap)
	d.drawMovingRectangle(d.screen)
	// 2 tegn world på screen. Camera applyes ved å justere op her.

	op := &ebiten.DrawImageOptions{}

	d.geomX -= 1
	op.GeoM.Translate(d.geomX, 0)

	screen.DrawImage(d.screen, op)
}

func (d *Drawer) DrawTiles(screen *ebiten.Image, worldMap *world_map.WorldMap) {
	// Get coordinates of the world map to draw. The viewport only shows a part of the world map, so we only draw that.
	// cameraViewport.Min.X is the x position of the viewport, similar for y.
	xCoordMin := d.mapViewport.Min.X
	xCoordMax := d.mapViewport.Max.X

	yCoordMin := d.mapViewport.Min.Y
	yCoordMax := d.mapViewport.Max.Y

	x := 0
	y := 0

	for yCoord := yCoordMin; yCoord < yCoordMax; yCoord++ {
		for xCoord := xCoordMin; xCoord <= xCoordMax; xCoord++ {
			tile := worldMap.TileAt(xCoord, yCoord)

			tileType := tile % 3
			var tileImage *ebiten.Image
			switch tileType {
			case Grass:
				tileImage = d.grassImage
			case Water:
				tileImage = d.waterImage
			case Mountain:
				tileImage = d.mountainImage
			default:
				panic("Unknown tile type")
			}

			// Draw tile at x, y
			op := &ebiten.DrawImageOptions{}
			//op.GeoM.Scale(scaleFactor, scaleFactor)
			// Set the image's pixel position
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(tileImage, op)

			x += TileSize
		}

		y += TileSize
		x = 0
	}
}

// TODO: CameraViewport must be at the very least have an x and y value which marks the top left corner of the screen.
// I.e. use pixels, not coordinates.
// We could use a rectangle to get max values as well, but screenWidth and screenHeight is the actual numbers we should use
// for drawing.
//func (d *Drawer) Draw(screen *ebiten.Image, worldMap *world_map.WorldMap, cameraViewport *image.Rectangle, scaleFactor float64) {
//	screenWidth := screen.Bounds().Dx()
//	screenHeight := screen.Bounds().Dy()
//
//	scaledTileSize := int(TileSize * scaleFactor)
//
//	for y := 0; y < screenHeight; y += scaledTileSize {
//		for x := 0; x < screenWidth; x += scaledTileSize {
//			xCoord := cameraViewport.Min.X
//			tile := worldMap.TileAt(xCoord, yCoord)
//		}
//	}
//}

func (d *Drawer) drawMovingRectangle(screen *ebiten.Image) {
	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	rectangleWidth := d.middleRectangle.Bounds().Dx()
	rectangleHeight := d.middleRectangle.Bounds().Dy()

	// Calculate the x and y coordinates to draw the image at the center of the window.
	x := float64(screenWidth/2-rectangleWidth/2) + d.rectangleX
	y := float64(screenHeight/2-rectangleHeight/2) + d.rectangleY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	screen.DrawImage(d.middleRectangle, op)
}

func (d *Drawer) Layout(outsideWidth, outsideHeight int, worldMap *world_map.WorldMap) {
	if outsideHeight != d.outsideHeight || outsideWidth != d.outsideWidth {
		log.Print("Window resized to ", outsideWidth, "x", outsideHeight)
		d.outsideWidth = outsideWidth
		d.outsideHeight = outsideHeight

		if !d.viewportInited {
			d.mapViewport = getMapViewportOfMapCenter(outsideWidth, outsideHeight, worldMap)
			d.screen = ebiten.NewImage(outsideWidth, outsideHeight)
			d.viewportInited = true
		}
	}
}

func (d *Drawer) IncreaseScaleFactor() {
	d.scaleFactor *= 1.01
}

func (d *Drawer) DecreaseScaleFactor() {
	d.scaleFactor *= 0.99
}

func (d *Drawer) MoveRectangle() {
	d.rectangleX += float64(rand.Intn(3)) - 1
	d.rectangleY += float64(rand.Intn(3)) - 1
}

func NewDrawer() *Drawer {
	d := &Drawer{}

	d.grassImage = ebiten.NewImage(TileSize, TileSize)
	d.grassImage.Fill(color.NRGBA{R: 188, G: 231, B: 132, A: 0xff})

	d.mountainImage = ebiten.NewImage(TileSize, TileSize)
	d.mountainImage.Fill(color.NRGBA{R: 66, G: 62, B: 55, A: 0xff})

	d.waterImage = ebiten.NewImage(TileSize, TileSize)
	d.waterImage.Fill(color.NRGBA{R: 52, G: 138, B: 167, A: 0xff})

	d.middleRectangle = ebiten.NewImage(50, 50)
	d.middleRectangle.Fill(color.NRGBA{R: 0x80, G: 0, B: 0, A: 0xff})

	d.scaleFactor = 1.0

	return d
}
