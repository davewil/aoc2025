package main

import (
	"fmt"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/aclements/go-z3/z3"
	"aoc2025/day10/parser"
	utils "github.com/davewil/aoc-utils"
)

func part1(lines parser.Input) int {
	total := 0
	for _, row := range lines.Rows {
		total += findMinPresses(row)
	}
	return total
}

func findMinPresses(row parser.Row) int {
	n := rowLength(row)
	if n == 0 {
		return 0
	}
	target := maskFromIndices(row.Lights)
	if target == 0 {
		return 0
	}
	buttonMasks := make([]int, len(row.Buttons))
	for i, b := range row.Buttons {
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

func part2(lines parser.Input) int {
	numWorkers := runtime.NumCPU()
	jobs := make(chan parser.Row, len(lines.Rows))
	results := make(chan int, len(lines.Rows))
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Go(func() {
			config := z3.NewContextConfig()
			ctx := z3.NewContext(config)
			solver := z3.NewSolver(ctx)

			// Hoist constants and allocate variable pool
			intSort := ctx.IntSort()
			zero := ctx.FromInt(0, intSort).(z3.Int)
			varPool := make([]z3.Int, 0, 100)

			for row := range jobs {
				solver.Reset()

				// Grow pool if needed
				for len(varPool) < len(row.Buttons) {
					varPool = append(varPool, ctx.IntConst(fmt.Sprintf("b%d", len(varPool))))
				}

				v, err := findMinPressesWithJoltagesILP(ctx, solver, row, varPool, zero, intSort)
				if err != nil {
					results <- -1
				} else {
					results <- v
				}
			}
		})
	}

	for _, row := range lines.Rows {
		jobs <- row
	}
	close(jobs)

	wg.Wait()
	close(results)

	total := 0
	for v := range results {
		if v == -1 {
			return -1
		}
		total += v
	}
	return total
}

// Uses Z3 SMT solver with a single context and incremental push/pop for binary search.
// Finds minimum sum of button presses subject to hitting exact joltages on each light.
func findMinPressesWithJoltagesILP(ctx *z3.Context, solver *z3.Solver, row parser.Row, varPool []z3.Int, zero z3.Int, intSort z3.Sort) (int, error) {
	if len(row.Joltages) == 0 {
		return 0, nil
	}
	if len(row.Buttons) == 0 {
		for _, v := range row.Joltages {
			if v != 0 {
				return -1, fmt.Errorf("no buttons to satisfy joltages")
			}
		}
		return 0, nil
	}

	nButtons := len(row.Buttons)
	nLights := len(row.Joltages)

	// Use pooled variables
	buttonVars := varPool[:nButtons]

	// Each button >= 0
	for i := range nButtons {
		solver.Assert(buttonVars[i].GE(zero))
	}

	// Sum of button contributions = joltage for each light
	for lightIdx := range nLights {
		var terms []z3.Int
		for bIdx, btn := range row.Buttons {
			if slices.Contains(btn, lightIdx) {
				terms = append(terms, buttonVars[bIdx])
			}
		}
		if len(terms) == 0 {
			if row.Joltages[lightIdx] != 0 {
				return -1, fmt.Errorf("light %d needs %d but no button affects it", lightIdx, row.Joltages[lightIdx])
			}
			continue
		}
		sum := terms[0].Add(terms[1:]...)
		solver.Assert(sum.Eq(ctx.FromInt(int64(row.Joltages[lightIdx]), intSort).(z3.Int)))
	}

	// Build total presses expression once
	totalPresses := buttonVars[0].Add(buttonVars[1:]...)

	// Upper bound
	maxSum := 0
	for _, j := range row.Joltages {
		maxSum += j
	}

	// Lower bound estimation
	// Sum(b_j * K_j) = TotalJoltage
	// S * MaxK >= TotalJoltage => S >= TotalJoltage / MaxK
	maxK := 0
	for _, btn := range row.Buttons {
		if len(btn) > maxK {
			maxK = len(btn)
		}
	}
	lo := 0
	if maxK > 0 {
		lo = (maxSum + maxK - 1) / maxK // ceil(maxSum / maxK)
	}

	// Model-guided binary search
	hi := maxSum
	best := -1

	for lo <= hi {
		mid := (lo + hi) / 2
		solver.Push()
		solver.Assert(totalPresses.LE(ctx.FromInt(int64(mid), intSort).(z3.Int)))
		sat, err := solver.Check()
		if err != nil {
			solver.Pop()
			return -1, err
		}
		if sat {
			// Found a solution, try to find a better one
			model := solver.Model()
			valExpr := model.Eval(totalPresses, true).(z3.Int)
			val, _, _ := valExpr.AsInt64()
			// model.Close() // Not available in go-z3

			best = int(val)
			hi = int(val) - 1
		} else {
			lo = mid + 1
		}
		solver.Pop()
	}

	if best == -1 {
		return -1, fmt.Errorf("no solution found")
	}
	return best, nil
}

func rowLength(r parser.Row) int {
	max := 0
	for _, idx := range r.Lights {
		if idx+1 > max {
			max = idx + 1
		}
	}
	for _, b := range r.Buttons {
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
	lines, err := parser.ParseInput(raw)
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
