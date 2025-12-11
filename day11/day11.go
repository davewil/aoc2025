package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	utils "github.com/davewil/aoc-utils"
)

type Graph struct {
	Adj   [][]int
	Nodes map[string]int
}

func parseInput(raw string) (Graph, error) {
	nodes := make(map[string]int)
	nextID := 0

	getID := func(name string) int {
		if id, ok := nodes[name]; ok {
			return id
		}
		id := nextID
		nextID++
		nodes[name] = id
		return id
	}

	// First pass: collect edges and assign IDs
	type edge struct {
		from, to int
	}
	// Pre-allocate edges slice to avoid resizing (guess size based on raw length)
	edges := make([]edge, 0, len(raw)/10)

	lines := strings.Split(strings.TrimSpace(raw), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}
		fromID := getID(parts[0])
		children := strings.Fields(parts[1])
		for _, child := range children {
			toID := getID(child)
			edges = append(edges, edge{fromID, toID})
		}
	}

	// Build adjacency list
	adj := make([][]int, nextID)
	for _, e := range edges {
		adj[e.from] = append(adj[e.from], e.to)
	}

	return Graph{Adj: adj, Nodes: nodes}, nil
}

func countPaths(adj [][]int, current, target int, memo []int) int {
	if current == target {
		return 1
	}
	if memo[current] != -1 {
		return memo[current]
	}

	total := 0
	for _, neighbor := range adj[current] {
		total += countPaths(adj, neighbor, target, memo)
	}
	memo[current] = total
	return total
}

func part1(graph Graph) int {
	start, ok1 := graph.Nodes["you"]
	end, ok2 := graph.Nodes["out"]
	if !ok1 || !ok2 {
		return 0
	}

	memo := make([]int, len(graph.Adj))
	for i := range memo {
		memo[i] = -1
	}
	return countPaths(graph.Adj, start, end, memo)
}

func part2(graph Graph) int {
	svr, ok1 := graph.Nodes["svr"]
	dac, ok2 := graph.Nodes["dac"]
	fft, ok3 := graph.Nodes["fft"]
	out, ok4 := graph.Nodes["out"]

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return 0
	}

	// Reuse single memo buffer
	memo := make([]int, len(graph.Adj))

	// Helper to run countPaths with reused memo
	runCount := func(from, to int) int {
		for i := range memo {
			memo[i] = -1
		}
		return countPaths(graph.Adj, from, to, memo)
	}

	// Check connectivity between dac and fft to determine order
	dacToFft := runCount(dac, fft)
	fftToDac := runCount(fft, dac)

	if dacToFft > 0 {
		// Order: svr -> dac -> fft -> out
		p1 := runCount(svr, dac)
		p2 := dacToFft // Already calculated
		p3 := runCount(fft, out)
		return p1 * p2 * p3
	} else if fftToDac > 0 {
		// Order: svr -> fft -> dac -> out
		p1 := runCount(svr, fft)
		p2 := fftToDac // Already calculated
		p3 := runCount(dac, out)
		return p1 * p2 * p3
	}

	return 0
}

func toDOT(graph Graph) string {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")

	// Create reverse mapping for names
	idToName := make([]string, len(graph.Nodes))
	for name, id := range graph.Nodes {
		idToName[id] = name
	}

	// Highlight special nodes
	specialNodes := []string{"svr", "fft", "dac", "out"}
	for _, name := range specialNodes {
		if _, ok := graph.Nodes[name]; ok {
			sb.WriteString(fmt.Sprintf("  %s [style=filled, fillcolor=lightblue];\n", name))
		}
	}

	for from, neighbors := range graph.Adj {
		fromName := idToName[from]
		for _, to := range neighbors {
			toName := idToName[to]
			sb.WriteString(fmt.Sprintf("  %s -> %s;\n", fromName, toName))
		}
	}

	sb.WriteString("}")
	return sb.String()
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

	// Save Graphviz DOT output to file
	if err := os.WriteFile("graph.dot", []byte(toDOT(graph)), 0644); err != nil {
		fmt.Println("Error writing graph.dot:", err)
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
