package world_map

import (
	"math"
	"math/rand"
)

type WorldMap struct {
	data [][]int
}

func (m *WorldMap) TileAt(xCoord int, yCoord int) int {
	return m.data[xCoord][yCoord]
}

func (m *WorldMap) Height() int {
	return len(m.data)
}

func (m *WorldMap) Width() int {
	if len(m.data) == 0 {
		return 0
	}

	return len(m.data[0])
}

type Point struct {
	x, y int
}

func Generate(width, height, nSites int) *WorldMap {
	// Generate random sites.
	sites := make([]Point, nSites)
	for i := range sites {
		sites[i] = Point{rand.Intn(width), rand.Intn(height)}
	}

	worldMap := make([][]int, height)
	for i := range worldMap {
		worldMap[i] = make([]int, width)
	}

	// Assign each point to the nearest site.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nearest := nearestSite(x, y, sites)
			worldMap[y][x] = nearest
		}
	}

	// Print the worldMap with site indices.
	//for y := 0; y < height; y++ {
	//	fmt.Printf("y=%d: ", y)
	//
	//	for x := 0; x < width; x++ {
	//		fmt.Printf("%2d ", worldMap[y][x])
	//	}
	//	fmt.Println()
	//}

	return &WorldMap{
		data: worldMap,
	}
}

func nearestSite(x, y int, sites []Point) int {
	minDistance := math.MaxInt64
	nearest := -1
	for i, site := range sites {
		distance := manhattanDistance(x, y, site.x, site.y)
		if distance < minDistance {
			minDistance = distance
			nearest = i
		}
	}
	return nearest
}

func manhattanDistance(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
