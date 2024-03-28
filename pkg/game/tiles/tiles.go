package tiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Drawer struct {
	grassImage    *ebiten.Image
	mountainImage *ebiten.Image
	waterImage    *ebiten.Image
}

// Tile types.
const (
	Grass = iota
	Water
	Mountain
)

// Tile size.
const (
	tileSize = 32
)

// Map data.
//var worldMap = [][]int{
//	{Grass, Water, Mountain, Water, Grass, Grass, Grass},
//	{Water, Grass, Grass, Water, Grass, Grass, Grass},
//	{Mountain, Grass, Grass, Grass, Grass, Grass, Grass},
//	{Water, Grass, Mountain, Mountain, Grass, Water, Grass},
//}

func (td *Drawer) Draw(screen *ebiten.Image, worldMap [][]int) {
	// Iterate over the map and draw the tiles.
	for mapY, row := range worldMap {
		for mapX, tile := range row {
			tileType := tile % 3
			switch tileType {
			case Grass:
				td.drawAt(td.grassImage, mapX, mapY, screen)
			case Water:
				td.drawAt(td.waterImage, mapX, mapY, screen)
			case Mountain:
				td.drawAt(td.mountainImage, mapX, mapY, screen)
			}

			// Draw the tile image at (x * tileSize, y * tileSize).
		}
	}
}

func (td *Drawer) drawAt(img *ebiten.Image, mapX int, mapY int, screen *ebiten.Image) {
	x := mapX * tileSize
	y := mapY * tileSize

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(img, op)
}

func (td *Drawer) init() {
	td.grassImage = ebiten.NewImage(tileSize, tileSize)
	td.grassImage.Fill(color.NRGBA{R: 0x00, G: 255, B: 0, A: 0xff})

	td.mountainImage = ebiten.NewImage(tileSize, tileSize)
	td.mountainImage.Fill(color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xff})

	td.waterImage = ebiten.NewImage(tileSize, tileSize)
	td.waterImage.Fill(color.NRGBA{R: 0, G: 0, B: 200, A: 0xff})
}

func NewDrawer() *Drawer {
	td := &Drawer{}
	td.init()

	return td
}
