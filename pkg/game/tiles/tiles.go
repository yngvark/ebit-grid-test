package tiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

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

//Map data.

var worldMap = [][]int{
	{Grass, Water, Mountain, Water, Grass, Grass, Grass},
	{Water, Grass, Grass, Water, Grass, Grass, Grass},
	{Mountain, Grass, Grass, Grass, Grass, Grass, Grass},
	{Water, Grass, Mountain, Mountain, Grass, Water, Grass},
}

//var worldMap = world_map.Generate()

// Images
var grassImage *ebiten.Image
var mountainImage *ebiten.Image
var waterImage *ebiten.Image

func WorldMapWidth() int {
	return len(worldMap[0])
}

func WorldMapHeight() int {
	return len(worldMap)
}

func Init() {
	grassImage = ebiten.NewImage(tileSize, tileSize)
	grassImage.Fill(color.NRGBA{R: 0x00, G: 255, B: 0, A: 0xff})

	mountainImage = ebiten.NewImage(tileSize, tileSize)
	mountainImage.Fill(color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xff})

	waterImage = ebiten.NewImage(tileSize, tileSize)
	waterImage.Fill(color.NRGBA{R: 0, G: 0, B: 200, A: 0xff})
}

func Draw(screen *ebiten.Image) {

	// Iterate over the map and draw the tiles.
	for mapY, row := range worldMap {
		for mapX, tile := range row {
			switch tile {
			case Grass:
				// Draw grass tile.
				drawAt(grassImage, mapX, mapY, screen)
			case Water:
				// Draw water tile.
				drawAt(waterImage, mapX, mapY, screen)
			case Mountain:
				// Draw mountain tile.
				drawAt(mountainImage, mapX, mapY, screen)
			}

			// Draw the tile image at (x * tileSize, y * tileSize).
		}
	}
}

func drawAt(img *ebiten.Image, mapX int, mapY int, screen *ebiten.Image) {
	x := mapX * tileSize
	y := mapY * tileSize

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(img, op)
}
