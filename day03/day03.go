package main

import (
	"fmt"
	"slices"
	"time"

	utils "github.com/davewil/aoc-utils"
)

func parseInput(raw string) ([][]int, error) {
	lines := utils.ReadLines(raw)
	banks := make([][]int, len(lines))
	for i, line := range lines {
		nums := make([]int, len(line))
		for n := 0; n < len(line); n++ {
			num := int(line[n] - '0')
			nums[n] = num
		}
		banks[i] = nums
	}

	return banks, nil
}

func part1(banks [][]int) int {
	total := 0

	for _, bank := range banks {
		total += maxJoltage(bank, 2)
	}

	return total
}

func part2(banks [][]int, target int) int {
	var total int

	for _, bank := range banks {
		joltage := maxJoltage(bank, target)
		total += joltage
	}

	return total
}

func maxJoltage(bank []int, digits int) int {
	if digits == 1 {
		return slices.Max(bank)
	}

	searchRange := bank[:len(bank)-(digits-1)]
	maxDigit := slices.Max(searchRange)

	idx := slices.Index(bank, maxDigit)
	return maxDigit*pow10(digits-1) + maxJoltage(bank[idx+1:], digits-1)
}

func pow10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 3)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	input, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	start := time.Now()
	part1Result := part1(input)
	part1Duration := time.Since(start)
	fmt.Printf("Part 1: %d (took %v)\n", part1Result, part1Duration)

	start = time.Now()
	part2Result := part2(input, 12)
	part2Duration := time.Since(start)
	fmt.Printf("Part 2: %d (took %v)\n", part2Result, part2Duration)
}
