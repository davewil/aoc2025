package main

import (
	"fmt"
	"time"

	"aoc2025/utils"

	aocutils "github.com/davewil/aoc-utils"
)

const (
	maxNeighboursToRemove = 3
	targetRune            = '@'
)

var directions = []utils.Point{
	{X: 0, Y: -1},  // Up
	{X: 1, Y: -1},  // Up-Right
	{X: 1, Y: 0},   // Right
	{X: 1, Y: 1},   // Down-Right
	{X: 0, Y: 1},   // Down
	{X: -1, Y: 1},  // Down-Left
	{X: -1, Y: 0},  // Left
	{X: -1, Y: -1}, // Up-Left
}

type Accumulator struct {
	CurrentPosition utils.Point
	Removed         []utils.Point
	RemovedCount    int
}

func parseInput(raw string) (*aocutils.Grid[rune], error) {
	lines := aocutils.ReadLines(raw)
	grid, err := aocutils.NewGridFromLines[rune](lines)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

func part1(grid *aocutils.Grid[rune]) int {
	acc := Accumulator{CurrentPosition: utils.Point{X: 0, Y: 0}, Removed: []utils.Point{}, RemovedCount: 0}
	removedRolls := removeRolls(acc, grid, targetRune)
	return len(removedRolls)
}

func part2(grid *aocutils.Grid[rune]) int {
	acc := Accumulator{CurrentPosition: utils.Point{X: 0, Y: 0}, Removed: []utils.Point{}, RemovedCount: 0}
	removedRolls := keepRemovingRolls(acc, grid)
	return removedRolls
}

func removeRolls(acc Accumulator, grid *aocutils.Grid[rune], target rune) []utils.Point {
	if acc.CurrentPosition.Y >= grid.Rows {
		return acc.Removed
	}
	rollsAroundCurrentPosition := countNeighbours(grid, acc.CurrentPosition, target)
	if grid.Get(acc.CurrentPosition.X, acc.CurrentPosition.Y) == target && rollsAroundCurrentPosition <= maxNeighboursToRemove {
		acc.Removed = append(acc.Removed, acc.CurrentPosition)
	}
	nextPos, ok := nextPosition(acc.CurrentPosition, grid)
	if !ok {
		return acc.Removed
	}
	acc.CurrentPosition = nextPos
	return removeRolls(acc, grid, target)
}

func countNeighbours(grid *aocutils.Grid[rune], pos utils.Point, target rune) int {
	count := 0
	for _, dir := range directions {
		x, y := pos.X+dir.X, pos.Y+dir.Y
		if grid.InBounds(x, y) && grid.Get(x, y) == target {
			count++
		}
	}
	return count
}

func nextPosition(pos utils.Point, grid *aocutils.Grid[rune]) (utils.Point, bool) {
	next := utils.Point{X: pos.X + 1, Y: pos.Y}
	if next.X >= grid.Cols {
		next.X = 0
		next.Y++
	}
	if next.Y >= grid.Rows {
		return utils.Point{}, false
	}
	return next, true
}

func keepRemovingRolls(acc Accumulator, grid *aocutils.Grid[rune]) int {
	removedRolls := removeRolls(Accumulator{CurrentPosition: utils.Point{X: 0, Y: 0}}, grid, targetRune)
	if len(removedRolls) == 0 {
		return acc.RemovedCount
	}

	for _, coord := range removedRolls {
		grid.Set(coord.X, coord.Y, '.')
	}

	acc.RemovedCount += len(removedRolls)
	return keepRemovingRolls(acc, grid)
}

func main() {
	aocutils.LoadEnv()
	raw, err := aocutils.GetPuzzleInput(2025, 4)
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
	part1Duration := time.Since(start)
	fmt.Printf("Part 1: %d (took %v)\n", part1Result, part1Duration)

	start = time.Now()
	part2Result := part2(lines)
	part2Duration := time.Since(start)
	fmt.Printf("Part 2: %d (took %v)\n", part2Result, part2Duration)

	// Reset grid for part 3 since part 2 mutates it
	lines, _ = parseInput(raw)
	start = time.Now()
	part3Result := part3(lines)
	part3Duration := time.Since(start)
	fmt.Printf("Part 3: %d (took %v)\n", part3Result, part3Duration)

	// Reset grid for part 4
	lines, _ = parseInput(raw)
	start = time.Now()
	part4Result := part4(lines)
	part4Duration := time.Since(start)
	fmt.Printf("Part 4: %d (took %v)\n", part4Result, part4Duration)

	// Reset grid for part 5
	lines, _ = parseInput(raw)
	start = time.Now()
	part5Result := part5(lines)
	part5Duration := time.Since(start)
	fmt.Printf("Part 5: %d (took %v)\n", part5Result, part5Duration)
}

// part3 is the iterative version of part2.
// Benchmarks show it is ~3x faster than the recursive version (33ms vs 105ms).
func part3(grid *aocutils.Grid[rune]) int {
	totalRemoved := 0
	for {
		removedInRound := []utils.Point{}
		for y := 0; y < grid.Rows; y++ {
			for x := 0; x < grid.Cols; x++ {
				pos := utils.Point{X: x, Y: y}
				if grid.Get(x, y) == targetRune {
					count := countNeighbours(grid, pos, targetRune)
					if count <= maxNeighboursToRemove {
						removedInRound = append(removedInRound, pos)
					}
				}
			}
		}

		if len(removedInRound) == 0 {
			break
		}

		for _, coord := range removedInRound {
			grid.Set(coord.X, coord.Y, '.')
		}
		totalRemoved += len(removedInRound)
	}
	return totalRemoved
}

// part4 uses a map instead of a grid.
func part4(grid *aocutils.Grid[rune]) int {
	// Convert grid to map
	m := make(map[utils.Point]rune)
	for y := 0; y < grid.Rows; y++ {
		for x := 0; x < grid.Cols; x++ {
			m[utils.Point{X: x, Y: y}] = grid.Get(x, y)
		}
	}

	totalRemoved := 0
	rows, cols := grid.Rows, grid.Cols

	for {
		removedInRound := []utils.Point{}
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				pos := utils.Point{X: x, Y: y}
				if val, exists := m[pos]; exists && val == targetRune {
					count := 0
					for _, dir := range directions {
						neighbour := utils.Point{X: pos.X + dir.X, Y: pos.Y + dir.Y}
						if nVal, ok := m[neighbour]; ok && nVal == targetRune {
							count++
						}
					}
					if count <= maxNeighboursToRemove {
						removedInRound = append(removedInRound, pos)
					}
				}
			}
		}

		if len(removedInRound) == 0 {
			break
		}

		for _, coord := range removedInRound {
			m[coord] = '.'
		}
		totalRemoved += len(removedInRound)
	}
	return totalRemoved
}

// part5 uses a flattened 1D slice, direct access, and buffer reuse.
func part5(grid *aocutils.Grid[rune]) int {
	rows, cols := grid.Rows, grid.Cols
	// Flatten grid
	flat := make([]rune, rows*cols)
	for y := range rows {
		for x := range cols {
			flat[y*cols+x] = grid.Get(x, y)
		}
	}

	totalRemoved := 0
	// Pre-allocate buffer for removed indices (max possible is all cells)
	removedInRound := make([]int, 0, rows*cols)

	for {
		removedInRound = removedInRound[:0] // Reset buffer

		for y := range rows {
			rowOffset := y * cols
			for x := range cols {
				idx := rowOffset + x
				if flat[idx] == targetRune {
					count := 0
					for _, dir := range directions {
						nx, ny := x+dir.X, y+dir.Y
						if nx >= 0 && nx < cols && ny >= 0 && ny < rows {
							if flat[ny*cols+nx] == targetRune {
								count++
							}
						}
					}
					if count <= maxNeighboursToRemove {
						removedInRound = append(removedInRound, idx)
					}
				}
			}
		}

		if len(removedInRound) == 0 {
			break
		}

		for _, idx := range removedInRound {
			flat[idx] = '.'
		}
		totalRemoved += len(removedInRound)
	}
	return totalRemoved
}
