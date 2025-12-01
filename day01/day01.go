package main

import (
	"fmt"
	"regexp"
	"strconv"

	utils "github.com/davewil/aoc-utils"
)

type Instruction struct {
	Direction string
	Distance  int
}

var (
	instructionPattern = regexp.MustCompile(`(L|R)(\d+)`)
)

func parseInput(input string) ([]Instruction, error) {
	matches := instructionPattern.FindAllStringSubmatch(input, -1)
	fmt.Printf("Found %d instructions\n", len(matches))
	instructions := make([]Instruction, len(matches))
	for i, m := range matches {
		instructions[i].Direction = m[1]
		d, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		instructions[i].Distance = d
	}
	return instructions, nil
}

func part1(instructions []Instruction) int {
	currentNotch := 50
	notchZeroCount := 0

	for _, instr := range instructions {
		fmt.Printf("At notch %d, instruction: %s%d\n", currentNotch, instr.Direction, instr.Distance)
		switch instr.Direction {
		case "L":
			nextNotch := currentNotch - instr.Distance
			nextNotch %= 100
			if nextNotch < 0 {
				nextNotch += 100
			}
			currentNotch = nextNotch
		case "R":
			nextNotch := (currentNotch + instr.Distance) % 100
			currentNotch = nextNotch
		}
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
