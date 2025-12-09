package main

import (
	"fmt"
	"strconv"
	"strings"

	utils "github.com/davewil/aoc-utils"
)

type Point2D struct {
	X, Y int
}

func parseInput(raw string) ([]Point2D, error) {
	lines := utils.ReadLines(raw)
	coords := []Point2D{}
	for _, line := range lines {
		before, after, found := strings.Cut(line, ",")
		if found {
			x, err := strconv.Atoi(before)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %w", before, err)
			}
			y, err := strconv.Atoi(after)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %w", after, err)
			}
			coords = append(coords, Point2D{x, y})
		}
	}
	return coords, nil
}

func part1(coords []Point2D) int {
	maxArea := 0

	for from := 0; from < len(coords)-1; from++ {
		for to := from + 1; to < len(coords); to++ {
			fromPoint := coords[from]
			toPoint := coords[to]
			area := (utils.Abs(fromPoint.X-toPoint.X) + 1) * (utils.Abs(fromPoint.Y-toPoint.Y) + 1)
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}

type Edge struct {
	from, to Point2D
}

func buildPerimeterAndEdges(redTiles []Point2D) ([]Point2D, []Edge) {
	perimeter := make([]Point2D, 0, len(redTiles)*2)
	edges := make([]Edge, 0, len(redTiles))

	appendEdgePoints := func(fromPoint, toPoint Point2D) {
		if fromPoint.X == toPoint.X {
			// vertical edge
			step := 1
			if toPoint.Y < fromPoint.Y {
				step = -1
			}
			for y := fromPoint.Y + step; y != toPoint.Y; y += step {
				perimeter = append(perimeter, Point2D{X: fromPoint.X, Y: y})
			}
		} else if fromPoint.Y == toPoint.Y {
			// horizontal edge
			step := 1
			if toPoint.X < fromPoint.X {
				step = -1
			}
			for x := fromPoint.X + step; x != toPoint.X; x += step {
				perimeter = append(perimeter, Point2D{X: x, Y: fromPoint.Y})
			}
		}
	}

	for idx := range redTiles {
		fromPoint := redTiles[idx]
		toPoint := redTiles[(idx+1)%len(redTiles)]
		edges = append(edges, Edge{from: fromPoint, to: toPoint})
		perimeter = append(perimeter, fromPoint)
		appendEdgePoints(fromPoint, toPoint)
	}

	return perimeter, edges
}

func part2(redTiles []Point2D) int {
	if len(redTiles) == 0 {
		return 0
	}
	perimeter, edges := buildPerimeterAndEdges(redTiles)

	maxArea := 0
	for from := 0; from < len(redTiles)-1; from++ {
		for to := from + 1; to < len(redTiles); to++ {
			fromPoint := redTiles[from]
			toPoint := redTiles[to]

			minX, maxX := sort(fromPoint.X, toPoint.X)
			minY, maxY := sort(fromPoint.Y, toPoint.Y)
			if minX == maxX || minY == maxY {
				continue
			}

			if intersects(edges, fromPoint, toPoint) {
				continue
			}

			centerX := float64(minX+maxX) / 2
			centerY := float64(minY+maxY) / 2
			if !pointInPolygon(centerX, centerY, redTiles) {
				continue
			}

			if hasPerimeterPointInside(perimeter, minX, maxX, minY, maxY) {
				continue
			}

			area := (maxX - minX + 1) * (maxY - minY + 1)
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}

func intersects(edges []Edge, from Point2D, to Point2D) bool {
	minX, maxX := sort(from.X, to.X)
	minY, maxY := sort(from.Y, to.Y)

	for _, edge := range edges {
		if edge.from.X == edge.to.X {
			x := edge.from.X
			if x <= minX || x >= maxX {
				continue
			}
			edgeMinY, edgeMaxY := sort(edge.from.Y, edge.to.Y)
			if strictOverlap(minY, maxY, edgeMinY, edgeMaxY) {
				return true
			}
		} else {
			y := edge.from.Y
			if y <= minY || y >= maxY {
				continue
			}
			edgeMinX, edgeMaxX := sort(edge.from.X, edge.to.X)
			if strictOverlap(minX, maxX, edgeMinX, edgeMaxX) {
				return true
			}
		}
	}
	return false
}

func sort(i1, i2 int) (int, int) {
	if i1 < i2 {
		return i1, i2
	}
	return i2, i1
}

func strictOverlap(aMin, aMax, bMin, bMax int) bool {
	return aMax > bMin && bMax > aMin
}

func hasPerimeterPointInside(perimeter []Point2D, minX, maxX, minY, maxY int) bool {
	for _, p := range perimeter {
		if p.X > minX && p.X < maxX && p.Y > minY && p.Y < maxY {
			return true
		}
	}
	return false
}

func pointInPolygon(x, y float64, polygon []Point2D) bool {
	inside := false
	for i := range polygon {
		j := (i + len(polygon) - 1) % len(polygon)
		xi, yi := float64(polygon[i].X), float64(polygon[i].Y)
		xj, yj := float64(polygon[j].X), float64(polygon[j].Y)
		intersect := ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
	}
	return inside
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 9)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	lines, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
}
