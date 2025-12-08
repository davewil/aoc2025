package main

import (
	"aoc2025/utils"
	"fmt"

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
func part2(lines *aocutils.Grid[rune]) int {
	return 0
}

func main() {
	aocutils.LoadEnv()
	raw, err := aocutils.GetPuzzleInput(2025, 7)
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
