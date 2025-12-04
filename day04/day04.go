package main

import (
	"fmt"

	utils "github.com/davewil/aoc-utils"
)

type Coord2D struct {
	X, Y int
}

var directions = [][2]int{
	{0, -1},  // Up
	{1, -1},  // Up-Right
	{1, 0},   // Right
	{1, 1},   // Down-Right
	{0, 1},   // Down
	{-1, 1},  // Down-Left
	{-1, 0},  // Left
	{-1, -1}, // Up-Left
}

type Accumulator struct {
	CurrentPosition Coord2D
	Removed         []Coord2D
	RemovedCount    int
}

func nextPosition(pos Coord2D, grid *utils.Grid[rune]) (Coord2D, bool) {
	next := Coord2D{X: pos.X + 1, Y: pos.Y}
	if next.X >= grid.Cols {
		next.X = 0
		next.Y++
	}
	if next.Y >= grid.Rows {
		return Coord2D{}, false
	}
	return next, true
}

func parseInput(raw string) (*utils.Grid[rune], error) {
	lines := utils.ReadLines(raw)
	grid, err := utils.NewGridFromLines[rune](lines)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

func part1(grid *utils.Grid[rune]) int {
	acc := Accumulator{CurrentPosition: Coord2D{X: 0, Y: 0}, Removed: []Coord2D{}, RemovedCount: 0}
	removedRolls, _ := removeRolls(acc, grid)
	return len(removedRolls)
}

func removeRolls(acc Accumulator, grid *utils.Grid[rune]) ([]Coord2D, int) {
	if acc.CurrentPosition.Y >= grid.Rows && len(acc.Removed) == 0 {
		return acc.Removed, 0
	}
	rollsAroundCurrentPosition := 0
	for _, dir := range directions {
		newX := acc.CurrentPosition.X + dir[0]
		newY := acc.CurrentPosition.Y + dir[1]
		if grid.InBounds(newX, newY) {
			if grid.Get(newX, newY) == '@' {
				rollsAroundCurrentPosition++
			}
		}
	}
	if grid.Get(acc.CurrentPosition.X, acc.CurrentPosition.Y) == '@' && rollsAroundCurrentPosition < 4 {
		acc.Removed = append(acc.Removed, acc.CurrentPosition)
	}
	nextPos, ok := nextPosition(acc.CurrentPosition, grid)
	if !ok {
		return acc.Removed, len(acc.Removed)
	}
	acc.CurrentPosition = nextPos
	acc.RemovedCount = len(acc.Removed)
	return removeRolls(acc, grid)
}

func part2(grid *utils.Grid[rune]) int {
	acc := Accumulator{CurrentPosition: Coord2D{X: 0, Y: 0}, Removed: []Coord2D{}, RemovedCount: 0}
	removedRolls := keepRemovingRolls(acc, grid)
	return removedRolls
}

func keepRemovingRolls(acc Accumulator, grid *utils.Grid[rune]) int {
	removedRolls, removedCount := removeRolls(acc, grid)
	if removedCount == 0 {
		return acc.RemovedCount
	}

	for _, coord := range removedRolls {
		grid.Set(coord.X, coord.Y, '.')
	}

	newAcc := Accumulator{CurrentPosition: Coord2D{X: 0, Y: 0}, Removed: []Coord2D{}, RemovedCount: acc.RemovedCount + removedCount}
	return keepRemovingRolls(newAcc, grid)
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 4)
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
