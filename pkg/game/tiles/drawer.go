package tiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles/world_map"
	"image"
	"image/color"
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
}

func (d *Drawer) Draw(screen *ebiten.Image, worldMap *world_map.WorldMap, cameraViewport *image.Rectangle, scaleFactor float64) {
	// Get coordinates of the world map to draw. The viewport only shows a part of the world map, so we only draw that.
	// cameraViewport.Min.X is the x position of the viewport, similar for y.
	xCoordMin := cameraViewport.Min.X
	xCoordMax := cameraViewport.Max.X

	yCoordMin := cameraViewport.Min.Y
	yCoordMax := cameraViewport.Max.Y

	x := 0
	y := 0

	for yCoord := yCoordMin; yCoord < yCoordMax; yCoord++ {
		for xCoord := xCoordMin; xCoord <= xCoordMax; xCoord++ {
			tile := worldMap.TileAt(xCoord, yCoord)
			tileType := tile % 3

			var image *ebiten.Image

			switch tileType {
			case Grass:
				image = d.grassImage
			case Water:
				image = d.waterImage
			case Mountain:
				image = d.mountainImage
			}

			d.drawAt(image, x, y, screen, scaleFactor)

			x += TileSize
		}

		y += TileSize
		x = 0
	}
}

// Draw the tile image at logical coordinates mapX and mapY
func (d *Drawer) drawAt(img *ebiten.Image, x int, y int, screen *ebiten.Image, scaleFactor float64) {
	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(scaleFactor, scaleFactor)

	// Set the image's pixel position
	op.GeoM.Translate(float64(x), float64(y))

	// Draw the tile
	screen.DrawImage(img, op)
}

func NewDrawer() *Drawer {
	d := &Drawer{}

	d.grassImage = ebiten.NewImage(TileSize, TileSize)
	d.grassImage.Fill(color.NRGBA{R: 188, G: 231, B: 132, A: 0xff})

	d.mountainImage = ebiten.NewImage(TileSize, TileSize)
	d.mountainImage.Fill(color.NRGBA{R: 66, G: 62, B: 55, A: 0xff})

	d.waterImage = ebiten.NewImage(TileSize, TileSize)
	d.waterImage.Fill(color.NRGBA{R: 52, G: 138, B: 167, A: 0xff})

	return d
}
