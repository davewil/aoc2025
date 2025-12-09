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

func part2(coords []Point2D) int {
	return 0
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
