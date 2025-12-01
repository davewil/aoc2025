package main

import (
	"fmt"

	utils "github.com/davewil/aoc-utils"
)

func parseInput(input string) ([]string, error) {
	lines := utils.ReadLines(input)
	return lines, nil
}

func part1(input []string) int {
	// Solve part 1 here
	return 0
}

func part2(input []string) int {
	// Solve part 2 here
	return 0
}

func main() {
	utils.LoadEnv()

	input, err := utils.GetPuzzleInput(2025, 1)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}

	parsed, err := parseInput(input)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	fmt.Println("Part 1:", part1(parsed))
	fmt.Println("Part 2:", part2(parsed))
}
