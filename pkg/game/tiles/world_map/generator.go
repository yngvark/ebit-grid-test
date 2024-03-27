package world_map

import (
	"math"
	"math/rand"
)

const (
	width  = 80
	height = 80
	nSites = 100 // Number of Voronoi sites
)

type Point struct {
	x, y int
}

func Generate() [][]int {
	//rand.Seed(time.Now().UnixNano())

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

	return worldMap
	//
	//// Print the worldMap with site indices.
	//for y := 0; y < height; y++ {
	//	for x := 0; x < width; x++ {
	//		fmt.Printf("%2d ", worldMap[y][x])
	//	}
	//	fmt.Println()
	//}
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
