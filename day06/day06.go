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
	rotated     *utils.Grid[rune]
}

func parseInput(raw string) (Input, error) {
	lines := utils.ReadLines(raw)
	data := [][]string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		data = append(data, fields)
	}

	grid, _ := utils.NewGridFromLines[rune](lines)
	rotated := grid.RotateLeft()

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
		rotated:     rotated,
	}
	return input, nil
}

func applyOp(op string, values []int) int {
	switch op {
	case "+":
		sum := 0
		for _, v := range values {
			sum += v
		}
		return sum
	case "*":
		prod := 1
		for _, v := range values {
			prod *= v
		}
		return prod
	}
	return 0
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
		total += applyOp(op, col)
	}
	return total
}

func part2(input Input) int {
	total := 0
	operands := make([]int, 0)
	operator := ""
	applyPending := func() {
		if len(operands) == 0 || operator == "" {
			return
		}
		total += applyOp(operator, operands)
		operands = operands[:0]
		operator = ""
	}
	for r := 0; r < input.rotated.Rows; r++ {
		line := string(input.rotated.Row(r))
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			applyPending()
			continue
		}
		if idx := strings.IndexAny(line, "+*"); idx != -1 {
			numberPart := strings.TrimSpace(line[:idx])
			if numberPart != "" {
				if num, err := strconv.Atoi(numberPart); err == nil {
					operands = append(operands, num)
				}
			}
			operator = string(line[idx])
			applyPending()
			continue
		}
		if num, err := strconv.Atoi(trimmed); err == nil {
			operands = append(operands, num)
		}
	}
	applyPending()

	return total
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
