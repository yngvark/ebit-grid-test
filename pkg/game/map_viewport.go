package game

import (
	"github.com/yngvark/ebit-grid-test/pkg/game/tiles/world_map"
	"image"
	"math"
)

// getMapViewportOfMapCenter calculates which part of the map should be visible inside the viewport.
func getMapViewportOfMapCenter(windowWidth int, windowHeight int, worldMap *world_map.WorldMap) *image.Rectangle {
	windowWidthInCoords := int(math.Ceil(float64(windowWidth) / TileSize))
	windowHeightInCoords := int(math.Ceil(float64(windowHeight) / TileSize))

	// Use map center for now. Later we can use a camera position.
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
