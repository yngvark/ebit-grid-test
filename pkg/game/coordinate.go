package game

import "fmt"

type Coordinate struct {
	X, Y int
}

func (c Coordinate) String() any {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{x, y}
}
