package main

import (
	"slices"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aclements/go-z3/z3"
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
		parts := strings.SplitSeq(b, ",")
		for part := range parts {
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
	parts := strings.SplitSeq(s, ",")
	for part := range parts {
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
	total := 0
	for _, row := range lines.rows {
		v, err := findMinPressesWithJoltagesILP(row)
		if err != nil {
			return -1
		}
		total += v
	}
	return total
}

// Uses Z3 SMT solver with a single context and incremental push/pop for binary search.
// Finds minimum sum of button presses subject to hitting exact joltages on each light.
func findMinPressesWithJoltagesILP(row Row) (int, error) {
	if len(row.joltages) == 0 {
		return 0, nil
	}
	if len(row.buttons) == 0 {
		for _, v := range row.joltages {
			if v != 0 {
				return -1, fmt.Errorf("no buttons to satisfy joltages")
			}
		}
		return 0, nil
	}

	nButtons := len(row.buttons)
	nLights := len(row.joltages)

	// Create single context and solver
	config := z3.NewContextConfig()
	ctx := z3.NewContext(config)
	solver := z3.NewSolver(ctx)

	// Create integer variables for each button press count
	buttonVars := make([]z3.Int, nButtons)
	for i := range nButtons {
		buttonVars[i] = ctx.IntConst(fmt.Sprintf("b%d", i))
	}

	// Each button >= 0
	zero := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
	for i := range nButtons {
		solver.Assert(buttonVars[i].GE(zero))
	}

	// Sum of button contributions = joltage for each light
	for lightIdx := range nLights {
		var terms []z3.Int
		for bIdx, btn := range row.buttons {
			if slices.Contains(btn, lightIdx) {
					terms = append(terms, buttonVars[bIdx])
				}
		}
		if len(terms) == 0 {
			if row.joltages[lightIdx] != 0 {
				return -1, fmt.Errorf("light %d needs %d but no button affects it", lightIdx, row.joltages[lightIdx])
			}
			continue
		}
		sum := terms[0]
		for _, t := range terms[1:] {
			sum = sum.Add(t)
		}
		solver.Assert(sum.Eq(ctx.FromInt(int64(row.joltages[lightIdx]), ctx.IntSort()).(z3.Int)))
	}

	// Build total presses expression once
	totalPresses := buttonVars[0]
	for _, bv := range buttonVars[1:] {
		totalPresses = totalPresses.Add(bv)
	}

	// Upper bound
	maxSum := 0
	for _, j := range row.joltages {
		maxSum += j
	}

	// Binary search with push/pop
	lo, hi := 0, maxSum
	for lo < hi {
		mid := (lo + hi) / 2
		solver.Push()
		solver.Assert(totalPresses.LE(ctx.FromInt(int64(mid), ctx.IntSort()).(z3.Int)))
		sat, _ := solver.Check()
		solver.Pop()
		if sat {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	// Verify final answer
	solver.Push()
	solver.Assert(totalPresses.LE(ctx.FromInt(int64(lo), ctx.IntSort()).(z3.Int)))
	sat, _ := solver.Check()
	solver.Pop()
	if !sat {
		return -1, fmt.Errorf("no solution found")
	}
	return lo, nil
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

	start := time.Now()
	p1 := part1(lines)
	fmt.Printf("Part 1: %d (took %v)\n", p1, time.Since(start))

	start = time.Now()
	p2 := part2(lines)
	fmt.Printf("Part 2: %d (took %v)\n", p2, time.Since(start))
}
