package main

import (
	"fmt"
	utils "github.com/davewil/aoc-utils"
)

// parseInput is a placeholder; real parsing will be added once
// the Day 2 puzzle format is known.
func parseInput(raw string) ([]string, error) {
	return []string{}, nil
}

// part1 stub – implement when puzzle released.
func part1(lines []string) int { return 0 }

// part2 stub – implement when puzzle released.
func part2(lines []string) int { return 0 }

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 2)
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
