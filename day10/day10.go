package main

import (
	"fmt"
	"strconv"
	"strings"

	utils "github.com/davewil/aoc-utils"
)

type Row struct {
	lights   []int
	buttons  [][]int
	joltages []int
}

type Input struct {
	rows []Row
}

func parseInput(raw string) (Input, error) {
	lines := utils.ReadLines(raw)
	input := Input{}
	for _, line := range lines {
		row := Row{}
		fields := strings.Fields(line)
		if len(fields) > 1 {
			row.lights = parseLights(fields[0])
			row.buttons = parseButtons(fields[1 : len(fields)-1])
			joltages, _ := parseJoltages(fields[len(fields)-1])
			row.joltages = joltages
		}
		input.rows = append(input.rows, row)
	}
	return input, nil
}

func parseLights(s string) []int {
	s = strings.Trim(s, "[]")
	result := []int{}
	for i, c := range s {
		if c == '#' {
			result = append(result, i)
		}
	}
	return result
}

func parseButtons(s []string) [][]int {
	result := [][]int{}
	buttons := []int{}
	for _, b := range s {
		b = strings.Trim(b, "()")
		parts := strings.Split(b, ",")
		for _, part := range parts {
			p, err := strconv.Atoi(part)
			if err != nil {
				fmt.Println("Failed to convert button part:", err)
				continue
			}
			buttons = append(buttons, p)
		}
		result = append(result, buttons)
		buttons = []int{}
	}
	return result
}

func parseJoltages(s string) ([]int, error) {
	result := []int{}
	s = strings.Trim(s, "{}")
	parts := strings.Split(s, ",")
	for _, part := range parts {
		p, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("Failed to convet part: %v", err)
		}
		result = append(result, p)
	}
	return result, nil

}

func part1(lines Input) int {
	total := 0
	for _, row := range lines.rows {
		total += findMinPresses(row)
	}
	return total
}

func findMinPresses(row Row) int {
	n := rowLength(row)
	if n == 0 {
		return 0
	}
	target := maskFromIndices(row.lights)
	if target == 0 {
		return 0
	}
	buttonMasks := make([]int, len(row.buttons))
	for i, b := range row.buttons {
		buttonMasks[i] = maskFromIndices(b)
	}
	maxStates := 1 << n
	dist := make([]int, maxStates)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{0}
	dist[0] = 0

	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		if s == target {
			return dist[s]
		}
		for _, bm := range buttonMasks {
			ns := s ^ bm
			if dist[ns] == -1 {
				dist[ns] = dist[s] + 1
				q = append(q, ns)
			}
		}
	}
	return -1
}

func part2(lines Input) int {
	return 0
}

func rowLength(r Row) int {
	max := 0
	for _, idx := range r.lights {
		if idx+1 > max {
			max = idx + 1
		}
	}
	for _, b := range r.buttons {
		for _, idx := range b {
			if idx+1 > max {
				max = idx + 1
			}
		}
	}
	return max
}

func maskFromIndices(idxs []int) int {
	m := 0
	for _, i := range idxs {
		m |= 1 << i
	}
	return m
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 10)
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
