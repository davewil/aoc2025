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

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func crossesZero(oldMod, newMod, delta int) bool {
	if oldMod == 0 || delta == 0 || delta%100 == 0 {
		return false
	}
	if delta > 0 {
		return newMod <= oldMod
	}
	return newMod > oldMod || newMod == 0
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
	notchZeroSeenCount := 0
	for _, instr := range instructions {
		delta := instr.Distance
		if instr.Direction == 'L' {
			delta = -delta
		}
		currentNotch += delta
		if mod(currentNotch, 100) == 0 {
			notchZeroSeenCount++
		}
	}
	return notchZeroSeenCount
}

func part2(instructions []Instruction) int {
	currentNotch := 50
	notchZeroSeenCount := 0
	for _, instr := range instructions {
		delta := instr.Distance
		if instr.Direction == 'L' {
			delta = -delta
		}

		notchZeroSeenCount += absInt(delta) / 100

		oldMod := mod(currentNotch, 100)
		currentNotch += delta
		newMod := mod(currentNotch, 100)
		if crossesZero(oldMod, newMod, delta) {
			notchZeroSeenCount++
		}
	}
	return notchZeroSeenCount
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
