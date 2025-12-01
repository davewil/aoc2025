package main

import (
	"fmt"
	"regexp"
	"strconv"

	utils "github.com/davewil/aoc-utils"
)

type Instruction struct {
	Direction byte
	Distance  int
}

var (
	instructionPattern = regexp.MustCompile(`(L|R)(\d+)`)
)

func mod(a, m int) int {
	return ((a % m) + m) % m
}

func parseInput(input string) ([]Instruction, error) {
	matches := instructionPattern.FindAllStringSubmatch(input, -1)
	instructions := make([]Instruction, len(matches))
	for i, m := range matches {
		d, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		instructions[i] = Instruction{
			Direction: m[1][0],
			Distance:  d,
		}
	}
	return instructions, nil
}

func part1(instructions []Instruction) int {
	currentNotch := 50
	notchZeroCount := 0

	for _, instr := range instructions {
		delta := instr.Distance
		if instr.Direction == 'L' {
			delta = -delta
		}
		currentNotch = mod(currentNotch+delta, 100)
		if currentNotch == 0 {
			notchZeroCount++
		}
	}
	return notchZeroCount
}

func part2(instructions []Instruction) int {
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
