package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	utils "github.com/davewil/aoc-utils"
)

type Input struct {
	operators   []string
	expressions [][]int
}

func parseInput(raw string) (Input, error) {
	lines := utils.ReadLines(raw)
	data := [][]string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		data = append(data, fields)
	}

	slices.Reverse(data)

	expressions := [][]int{}
	for _, line := range data[1:] {
		expression := []int{}
		for _, field := range line {
			num, err := strconv.Atoi(field)
			if err != nil {
				return Input{}, err
			}
			expression = append(expression, num)
		}
		expressions = append(expressions, expression)
	}
	input := Input{
		operators:   data[0],
		expressions: expressions,
	}
	return input, nil
}

func part1(input Input) int {
	total := 0
	for i, op := range input.operators {
		col := []int{}
		for _, row := range input.expressions {
			if i < len(row) {
				col = append(col, row[i])
			}
		}

		switch op {
		case "+":
			sum := 0
			for _, v := range col {
				sum += v
			}
			total += sum
		case "*":
			prod := 1
			for _, v := range col {
				prod *= v
			}
			total += prod
		default:
			fmt.Printf("Unknown operator: %s\n", op)
		}

	}
	return total
}

func part2(input Input) int {
	return 0
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 6)
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
