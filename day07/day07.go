package main

import (
	"aoc2025/utils"
	"fmt"
	"time"

	aocutils "github.com/davewil/aoc-utils"
)

func parseInput(raw string) (*aocutils.Grid[rune], error) {
	lines := aocutils.ReadLines(raw)
	grid, err := aocutils.NewGridFromLines[rune](lines)
	if err != nil {
		return nil, err
	}

	return grid, nil
}
func part1(grid *aocutils.Grid[rune]) int {
	var start utils.Point
	for col := 0; col < grid.Cols; col++ {
		if grid.Get(0, col) == 'S' {
			start = utils.Point{X: col, Y: 0}
			break
		}
	}
	if start.Y+1 >= grid.Rows {
		return 0
	}
	grid.Set(start.Y+1, start.X, '|')
	for row := start.Y + 1; row < grid.Rows; row++ {
		for col := 0; col < grid.Cols; col++ {
			if grid.Get(row, col) == '|' {
				if row+1 < grid.Rows {
					below := grid.Get(row+1, col)
					switch below {
					case '^':
						if col-1 >= 0 {
							grid.Set(row+1, col-1, '|')
						}
						if col+1 < grid.Cols {
							grid.Set(row+1, col+1, '|')
						}
					case '.':
						grid.Set(row+1, col, '|')
					}
				}
			}
		}
	}
	collected := 0
	for row := 1; row < grid.Rows; row++ {
		for col := 0; col < grid.Cols; col++ {
			if grid.Get(row, col) == '^' && grid.Get(row-1, col) == '|' {
				collected++
			}
		}
	}
	return collected
}

type Range struct {
	Start, End int
}

func part2(grid *aocutils.Grid[rune]) int {
	// Use double buffering: current row counts and next row counts
	current := make([]int, grid.Cols)
	next := make([]int, grid.Cols)

	// Find starting position
	startCol := 0
	for col := 0; col < grid.Cols; col++ {
		if grid.Get(0, col) == 'S' {
			startCol = col
			break
		}
	}

	// Initialize: one route at the starting column
	current[startCol] = 1
	colRange := Range{startCol, startCol}

	for row := 1; row < grid.Rows; row++ {
		// Clear next row
		for i := range next {
			next[i] = 0
		}

		newStart := colRange.End // Will shrink to find actual min
		newEnd := colRange.Start // Will grow to find actual max

		for col := colRange.Start; col <= colRange.End; col++ {
			if current[col] == 0 {
				continue
			}

			cell := grid.Get(row, col)
			switch cell {
			case '^':
				// Splitter: routes go left and right
				if col-1 >= 0 {
					next[col-1] += current[col]
					if col-1 < newStart {
						newStart = col - 1
					}
					if col-1 > newEnd {
						newEnd = col - 1
					}
				}
				if col+1 < grid.Cols {
					next[col+1] += current[col]
					if col+1 < newStart {
						newStart = col + 1
					}
					if col+1 > newEnd {
						newEnd = col + 1
					}
				}
			case '.', 'S':
				// Empty cell: routes continue straight down
				next[col] += current[col]
				if col < newStart {
					newStart = col
				}
				if col > newEnd {
					newEnd = col
				}
			}
		}

		// Update colRange if we found any routes
		if newStart <= newEnd {
			colRange = Range{newStart, newEnd}
		}

		// Swap current and next
		current, next = next, current
	}

	// Sum all routes that made it to the bottom
	total := 0
	for _, v := range current {
		total += v
	}

	return total
}

func main() {
	aocutils.LoadEnv()
	raw, err := aocutils.GetPuzzleInput(2025, 7)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	grid1, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	grid2, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	start := time.Now()
	part1Result := part1(grid1)
	part1Duration := time.Since(start)
	fmt.Printf("Part 1: %d (took %v)\n", part1Result, part1Duration)

	start = time.Now()
	part2Result := part2(grid2)
	part2Duration := time.Since(start)
	fmt.Printf("Part 2: %d (took %v)\n", part2Result, part2Duration)
}
