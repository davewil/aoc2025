package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	utils "github.com/davewil/aoc-utils"
)

// Range struct duplicated here for standalone execution
type Range struct {
	Start int
	End   int
}

// parseInput duplicated here for standalone execution
func parseInput(input string) ([]Range, error) {
	parts := strings.Split(input, ",")
	ranges := make([]Range, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		bounds := strings.Split(p, "-")
		if len(bounds) != 2 {
			return nil, fmt.Errorf("invalid range: %s", p)
		}
		start, err := strconv.Atoi(bounds[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(bounds[1])
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, Range{Start: start, End: end})
	}
	return ranges, nil
}

// getDivisors duplicated here
func getDivisors(n int) []int {
	var result []int
	for i := 1; i <= n/2; i++ {
		if n%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

// hasRepeatingPattern duplicated here
func hasRepeatingPattern(s string, chunkLen int) bool {
	if chunkLen <= 0 || len(s)%chunkLen != 0 {
		return false
	}
	pattern := s[:chunkLen]
	for i := chunkLen; i < len(s); i += chunkLen {
		if s[i:i+chunkLen] != pattern {
			return false
		}
	}
	return true
}

// Part 1 Concurrent: Using a Worker Pool
func part1Concurrent(ranges []Range) int {
	// Channel for jobs (ranges)
	jobs := make(chan Range, len(ranges))
	// Channel for results (partial sums)
	results := make(chan int, len(ranges))

	// Number of workers to spawn
	numWorkers := 8 // Could be runtime.NumCPU()

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				partialSum := 0
				for i := r.Start; i <= r.End; i++ {
					s := strconv.Itoa(i)
					if len(s)%2 != 0 {
						continue
					}
					mid := len(s) / 2
					if s[:mid] == s[mid:] {
						partialSum += i
					}
				}
				results <- partialSum
			}
		}()
	}

	// Send jobs
	for _, r := range ranges {
		jobs <- r
	}
	close(jobs)

	// Wait for workers in a separate goroutine to close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	total := 0
	for res := range results {
		total += res
	}

	return total
}

// Part 2 Concurrent: Worker Pool on Ranges
func part2Concurrent(ranges []Range) int {
	// Channel for jobs
	jobs := make(chan Range, len(ranges))
	// Channel for results: each worker sends a slice of valid numbers found in a range
	results := make(chan []int, len(ranges))

	numWorkers := 8
	var wg sync.WaitGroup

	// Workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				var validNums []int
				for i := r.Start; i <= r.End; i++ {
					s := strconv.Itoa(i)
					if len(s) < 2 {
						continue
					}
					// Optimization: we don't check 'seen' here because we don't share state.
					// We just find all valid numbers in this range.
					divs := getDivisors(len(s))
					for _, d := range divs {
						if hasRepeatingPattern(s, d) {
							validNums = append(validNums, i)
							break
						}
					}
				}
				results <- validNums
			}
		}()
	}

	// Send jobs
	for _, r := range ranges {
		jobs <- r
	}
	close(jobs)

	// Wait and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Reducer: Deduplicate and sum
	seen := make(map[int]bool)
	total := 0
	for batch := range results {
		for _, num := range batch {
			if !seen[num] {
				seen[num] = true
				total += num
			}
		}
	}

	return total
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 2)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	data, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	fmt.Println("--- Concurrent Version ---")

	start := time.Now()
	result1 := part1Concurrent(data)
	fmt.Printf("Part 1: %d (took %v)\n", result1, time.Since(start))

	start = time.Now()
	result2 := part2Concurrent(data)
	fmt.Printf("Part 2: %d (took %v)\n", result2, time.Since(start))
}
