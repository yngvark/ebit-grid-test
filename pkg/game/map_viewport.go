package game

import (
	"image"
	"math"
)

// getMapViewportFromCoordinate calculates which part of the map should be visible inside the viewport.
func getMapViewportFromCoordinate(screenWidth int, screenHeight int, coordinate Coordinate) *image.Rectangle {
	// TODO: We need some handling if coordinate is near the edge of the worldmap. We dont want to show white void.
	windowWidthInCoords := int(math.Ceil(float64(screenWidth) / TileSize))
	windowHeightInCoords := int(math.Ceil(float64(screenHeight) / TileSize))

	// Use map center for now. Later we can use a camera position.
	xCoordsMin := int(math.Floor(float64(float64(coordinate.X) - float64(windowWidthInCoords)/2)))
	yCoordsMin := int(math.Floor(float64(float64(coordinate.Y) - float64(windowHeightInCoords/2))))

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
