package main

import (
	"fmt"

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
		total += getHighestPair(bank)
	}

	return total
}

func getHighestPair(bank []int) int {
	indices := make(map[int][]int)
	for i, val := range bank {
		indices[val] = append(indices[val], i)
	}

	for x := 9; x >= 0; x-- {
		xIndices, ok := indices[x]
		if !ok {
			continue
		}

		earliestX := xIndices[0]

		for y := 9; y >= 0; y-- {
			yIndices, ok := indices[y]
			if !ok {
				continue
			}

			lastY := yIndices[len(yIndices)-1]

			if lastY > earliestX {
				return x*10 + y
			}
		}
	}

	return -1
}

func part2(banks [][]int) int64 {
	var total int64

	for _, bank := range banks {
		joltage := getHighestJoltage(bank)
		total += joltage
	}

	return total
}

func getHighestJoltage(bank []int) int64 {
	b := make([]int, len(bank))
	copy(b, bank)

	i := 0
	//loop until only 12 digits remain
	for len(b) > 12 {
		//remove last digit if we're at the end. Either all equal or ever decreasing
		if i >= len(b)-1 {
			lastIdx := len(b) - 1
			b = b[:lastIdx]
			i = 0
		} else if b[i] < b[i+1] {
			//remove current digit if next digit is larger
			b = append(b[:i], b[i+1:]...)
			//if we're not at the start, step back one to re-evaluate previous digit
			if i > 0 {
				i--
			}
		} else {
			//keep going
			i++
		}
	}

	var result int64
	for _, digit := range b {
		result = result*10 + int64(digit)
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
	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
