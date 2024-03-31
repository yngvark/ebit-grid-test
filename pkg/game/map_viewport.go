package game

import (
	"image"
	"math"
)

// getMapViewportFromCoordinate calculates which part of the map should be visible inside the viewport.
func getMapViewportFromCoordinate(screenWidth int, screenHeight int, coordinate Coordinate) *image.Rectangle {
	//log.Printf("getMapViewportFromCoordinate: screenWidth: %d, screenHeight: %d, coordinate: %v",
	//	screenWidth, screenHeight, coordinate)
	// TODO: We need some handling if coordinate is near the edge of the worldmap. We dont want to show white void.
	windowWidthInCoords := int(math.Ceil(float64(screenWidth) / TileSize))
	windowHeightInCoords := int(math.Ceil(float64(screenHeight) / TileSize))

	// extraTiles below is there to show tiles when the user scrolls the map, i.e. when the camera changes position. This
	// makes the camera show the next tile.
	xCoordsMin := int(math.Floor(float64(coordinate.X)-float64(windowWidthInCoords)/2)) - extraTiles
	yCoordsMin := int(math.Floor(float64(coordinate.Y)-float64(windowHeightInCoords/2))) - extraTiles

	xCoordsMax := int(math.Ceil(float64(xCoordsMin+windowWidthInCoords))) + extraTiles
	yCoordsMax := int(math.Ceil(float64(yCoordsMin+windowHeightInCoords))) + extraTiles

	viewPortCoords := image.Rect(
		xCoordsMin,
		yCoordsMin,
		xCoordsMax,
		yCoordsMax,
	)

	r := viewPortCoords
	return &r
}

const extraTiles = 0
