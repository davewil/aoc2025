package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	utils "github.com/davewil/aoc-utils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Point3D struct {
	X, Y, Z int
}

// pointNode wraps a Point3D so the coordinates travel with the Gonum node.
type pointNode struct {
	id    int64
	Point Point3D
}

func (n *pointNode) ID() int64 {
	return n.id
}

func parseInput(raw string) ([]Point3D, *simple.WeightedUndirectedGraph, error) {
	lines := utils.ReadLines(raw)
	points := make([]Point3D, 0, len(lines))
	graph := simple.NewWeightedUndirectedGraph(0, math.Inf(1))
	var nextID int64
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		coords := strings.Split(line, ",")
		if len(coords) != 3 {
			return nil, nil, fmt.Errorf("line %d: expected 3 values, got %d", idx+1, len(coords))
		}
		var values [3]int
		for i, coord := range coords {
			num, err := strconv.Atoi(strings.TrimSpace(coord))
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: invalid integer %q: %w", idx+1, coord, err)
			}
			values[i] = num
		}
		point := Point3D{X: values[0], Y: values[1], Z: values[2]}
		points = append(points, point)
		node := &pointNode{id: nextID, Point: point}
		nextID++
		graph.AddNode(node)
	}
	return points, graph, nil
}

func connectNearestPairs(g *simple.WeightedUndirectedGraph, count int) {
	if g == nil || count <= 0 {
		return
	}
	var nodes []*pointNode
	iter := g.Nodes()
	for iter.Next() {
		if n, ok := iter.Node().(*pointNode); ok {
			nodes = append(nodes, n)
		}
	}
	if len(nodes) < 2 {
		return
	}
	type edgeCandidate struct {
		u, v *pointNode
		w    float64
	}
	edges := make([]edgeCandidate, 0, len(nodes)*(len(nodes)-1)/2)
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			u := nodes[i]
			v := nodes[j]
			edges = append(edges, edgeCandidate{
				u: u,
				v: v,
				w: straightLineDistance(u.Point, v.Point),
			})
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})
	added := 0
	for _, edge := range edges {
		if added >= count {
			break
		}
		if g.HasEdgeBetween(edge.u.ID(), edge.v.ID()) {
			continue
		}
		g.SetWeightedEdge(simple.WeightedEdge{F: edge.u, T: edge.v, W: edge.w})
		added++
	}
}

func straightLineDistance(a, b Point3D) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func connectedComponentSizes(g graph.Graph) []int {
	if g == nil {
		return nil
	}
	visited := make(map[int64]bool)
	var sizes []int
	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if visited[n.ID()] {
			continue
		}
		size := 0
		stack := []int64{n.ID()}
		visited[n.ID()] = true
		for len(stack) > 0 {
			id := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			size++
			neighbors := g.From(id)
			for neighbors.Next() {
				neighbor := neighbors.Node()
				if visited[neighbor.ID()] {
					continue
				}
				visited[neighbor.ID()] = true
				stack = append(stack, neighbor.ID())
			}
		}
		sizes = append(sizes, size)
	}
	return sizes
}

func part1(points []Point3D, g *simple.WeightedUndirectedGraph, connectCount int) int {
	_ = points
	connectNearestPairs(g, connectCount)
	sizes := connectedComponentSizes(g)
	if len(sizes) == 0 {
		return 0
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	limit := min(len(sizes), 3)
	product := 1
	for i := range limit {
		product *= sizes[i]
	}
	return product
}
func part2(points []Point3D, g graph.Weighted) int {
	_, _ = points, g
	return 0
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 8)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	points, lavaGraph, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	start := time.Now()
	part1Result := part1(points, lavaGraph, 1000)
	fmt.Printf("Part 1: %d (took %s)\n", part1Result, time.Since(start))
	fmt.Println("Part 2:", part2(points, lavaGraph))
}
