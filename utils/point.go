package utils

// Point represents a 2D coordinate.
type Point struct {
	X, Y int
}

// ManhattanDistance calculates the Manhattan distance (|x1 - x2| + |y1 - y2|)
// between two points a and b.
func ManhattanDistance(a, b Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
