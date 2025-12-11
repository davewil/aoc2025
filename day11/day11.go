package main

import (
	"fmt"
	"strings"
	"time"

	utils "github.com/davewil/aoc-utils"
)

type Graph map[string][]string

func parseInput(raw string) (Graph, error) {
	graph := make(Graph)
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}
		node := parts[0]
		children := strings.Fields(parts[1])
		graph[node] = children
	}
	return graph, nil
}

func countPaths(graph Graph, current, target string, memo map[string]int) int {
	if current == target {
		return 1
	}
	if count, ok := memo[current]; ok {
		return count
	}

	total := 0
	for _, neighbor := range graph[current] {
		total += countPaths(graph, neighbor, target, memo)
	}
	memo[current] = total
	return total
}

func part1(graph Graph) int {
	memo := make(map[string]int)
	return countPaths(graph, "you", "out", memo)
}

func part2(graph Graph) int {
	// Check connectivity between dac and fft to determine order
	memo := make(map[string]int)
	dacToFft := countPaths(graph, "dac", "fft", memo)

	memo = make(map[string]int)
	fftToDac := countPaths(graph, "fft", "dac", memo)

	if dacToFft > 0 {
		// Order: svr -> dac -> fft -> out
		memo1 := make(map[string]int)
		p1 := countPaths(graph, "svr", "dac", memo1)

		memo2 := make(map[string]int)
		p2 := countPaths(graph, "dac", "fft", memo2)

		memo3 := make(map[string]int)
		p3 := countPaths(graph, "fft", "out", memo3)

		return p1 * p2 * p3
	} else if fftToDac > 0 {
		// Order: svr -> fft -> dac -> out
		memo1 := make(map[string]int)
		p1 := countPaths(graph, "svr", "fft", memo1)

		memo2 := make(map[string]int)
		p2 := countPaths(graph, "fft", "dac", memo2)

		memo3 := make(map[string]int)
		p3 := countPaths(graph, "dac", "out", memo3)

		return p1 * p2 * p3
	}

	return 0
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 11)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	graph, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	start := time.Now()
	part1Result := part1(graph)
	part1Duration := time.Since(start)
	fmt.Printf("Part 1: %d (took %v)\n", part1Result, part1Duration)

	start = time.Now()
	part2Result := part2(graph)
	part2Duration := time.Since(start)
	fmt.Printf("Part 2: %d (took %v)\n", part2Result, part2Duration)
}
