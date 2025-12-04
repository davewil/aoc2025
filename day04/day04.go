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
	acc := Accumulator{CurrentPosition: Coord2D{X: 0, Y: 0}, Removed: []Coord2D{}}
	removedRolls := removeRolls(acc, grid)
	return len(removedRolls)
}

func removeRolls(acc Accumulator, grid *utils.Grid[rune]) []Coord2D {
	if acc.CurrentPosition.Y >= grid.Rows {
		return acc.Removed
	}
	// Perform test if roll can be removed
	// For each direction, check if there is a roll in that direction if there are fewer than 4 rolls remove
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
		return acc.Removed
	}
	acc.CurrentPosition = nextPos
	return removeRolls(acc, grid)
}

func part2(grid *utils.Grid[rune]) int {
	return 0
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
