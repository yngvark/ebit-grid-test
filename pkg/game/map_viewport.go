package game

import (
	"github.com/rs/zerolog/log"
	"image"
	"math"
)

// getMapViewportFromCoordinate calculates which part of the map should be visible inside the viewport.
func getMapViewportFromCoordinate(screenWidth int, screenHeight int, coordinate Coordinate) *image.Rectangle {
	log.Printf("getMapViewportFromCoordinate: screenWidth: %d, screenHeight: %d, coordinate: %v",
		screenWidth, screenHeight, coordinate)
	// TODO: We need some handling if coordinate is near the edge of the worldmap. We dont want to show white void.
	windowWidthInCoords := int(math.Ceil(float64(screenWidth) / TileSize))
	windowHeightInCoords := int(math.Ceil(float64(screenHeight) / TileSize))

	// -1 and +1 below is there to show tiles when the user scrolls the map, i.e. when the camera changes position. This
	// makes the camera show the next tile.
	xCoordsMin := int(math.Floor(float64(coordinate.X)-float64(windowWidthInCoords)/2)) - 1
	yCoordsMin := int(math.Floor(float64(coordinate.Y)-float64(windowHeightInCoords/2))) - 1

	xCoordsMax := int(math.Ceil(float64(xCoordsMin+windowWidthInCoords))) + 1
	yCoordsMax := int(math.Ceil(float64(yCoordsMin+windowHeightInCoords))) + 1

	viewPortCoords := image.Rect(
		xCoordsMin,
		yCoordsMin,
		xCoordsMax,
		yCoordsMax,
	)

	r := viewPortCoords
	return &r
}
