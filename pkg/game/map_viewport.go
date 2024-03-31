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

	xCoordsMin := int(math.Floor(float64(coordinate.X) - float64(windowWidthInCoords)/2))
	yCoordsMin := int(math.Floor(float64(coordinate.Y) - float64(windowHeightInCoords/2)))

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
