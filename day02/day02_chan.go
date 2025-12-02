package main

import (
	"strconv"
	"sync"
)

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
