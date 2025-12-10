package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

type Rectangle struct {
	minX, minY, maxX, maxY int
}

func buildEdges(redTiles []Point2D) []Edge {
	edges := make([]Edge, 0, len(redTiles))
	for idx := range redTiles {
		fromPoint := redTiles[idx]
		toPoint := redTiles[(idx+1)%len(redTiles)]
		edges = append(edges, Edge{from: fromPoint, to: toPoint})
	}
	return edges
}

func part2(redTiles []Point2D) int {
	_, area := bestRectangle(redTiles)
	return area
}

func bestRectangle(redTiles []Point2D) (Rectangle, int) {
	best := Rectangle{}
	if len(redTiles) == 0 {
		return best, 0
	}
	best, area, _ := searchRectangles(redTiles, nil, 0)
	return best, area
}

// searchRectangles iterates all vertex pairs; when cb is non-nil it will be called
// on each new best or every sampleEvery steps (>0). Returns best rectangle, area, and steps executed.
func searchRectangles(redTiles []Point2D, cb func(Rectangle, int, int), sampleEvery int) (Rectangle, int, int) {
	best := Rectangle{}
	edges := buildEdges(redTiles)
	insideCache := make(map[int64]bool, len(redTiles)*len(redTiles))
	maxArea := 0
	steps := 0
	if sampleEvery <= 0 {
		sampleEvery = 0
	}
	for from := 0; from < len(redTiles)-1; from++ {
		for to := from + 1; to < len(redTiles); to++ {
			fromPoint := redTiles[from]
			toPoint := redTiles[to]
			steps++

			minX, maxX := sort(fromPoint.X, toPoint.X)
			minY, maxY := sort(fromPoint.Y, toPoint.Y)
			if minX == maxX || minY == maxY {
				continue
			}

			if intersects(edges, fromPoint, toPoint) {
				continue
			}

			centerSumX := minX + maxX
			centerSumY := minY + maxY
			cacheKey := (int64(centerSumX) << 32) | int64(uint32(centerSumY))
			inside, ok := insideCache[cacheKey]
			if !ok {
				inside = pointInPolygonCenter(centerSumX, centerSumY, redTiles)
				insideCache[cacheKey] = inside
			}
			if !inside {
				continue
			}

			if hasInteriorVertex(redTiles, minX, maxX, minY, maxY, fromPoint, toPoint) {
				continue
			}

			area := (maxX - minX + 1) * (maxY - minY + 1)
			if area > maxArea {
				maxArea = area
				best = Rectangle{minX: minX, minY: minY, maxX: maxX, maxY: maxY}
				if cb != nil {
					cb(best, maxArea, steps)
				}
			} else if cb != nil && sampleEvery > 0 && steps%sampleEvery == 0 {
				cb(best, maxArea, steps)
			}
		}
	}

	return best, maxArea, steps
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

func hasInteriorVertex(vertices []Point2D, minX, maxX, minY, maxY int, skipA, skipB Point2D) bool {
	for _, p := range vertices {
		if (p == skipA) || (p == skipB) {
			continue
		}
		if p.X > minX && p.X < maxX && p.Y > minY && p.Y < maxY {
			return true
		}
	}
	return false
}

func pointInPolygonCenter(x2, y2 int, polygon []Point2D) bool {
	inside := false
	for i := range polygon {
		j := (i + len(polygon) - 1) % len(polygon)
		yi2 := polygon[i].Y * 2
		yj2 := polygon[j].Y * 2
		if (yi2 > y2) == (yj2 > y2) {
			continue
		}
		xi2 := polygon[i].X * 2
		xj2 := polygon[j].X * 2
		dy := yj2 - yi2
		if dy == 0 {
			continue
		}
		dx := xj2 - xi2
		yDelta := y2 - yi2
		lhs := x2 * dy
		rhs := xi2*dy + dx*yDelta
		if (dy > 0 && lhs < rhs) || (dy < 0 && lhs > rhs) {
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
	start := time.Now()
	part1Result := part1(lines)
	fmt.Printf("Part 1: %d (took %s)\n", part1Result, time.Since(start))
	start = time.Now()
	part2Result := part2(lines)
	fmt.Printf("Part 2: %d (took %s)\n", part2Result, time.Since(start))
}
