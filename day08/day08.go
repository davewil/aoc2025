package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

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
	for range count {
		bestDist := math.Inf(1)
		var bestU, bestV *pointNode
		for i := 0; i < len(nodes); i++ {
			for j := i + 1; j < len(nodes); j++ {
				u := nodes[i]
				v := nodes[j]
				if g.HasEdgeBetween(u.ID(), v.ID()) {
					continue
				}
				d := straightLineDistance(u.Point, v.Point)
				if d < bestDist {
					bestDist = d
					bestU, bestV = u, v
				}
			}
		}
		if bestU == nil || bestV == nil {
			break
		}
		g.SetWeightedEdge(simple.WeightedEdge{F: bestU, T: bestV, W: bestDist})
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
	limit := 3
	if len(sizes) < limit {
		limit = len(sizes)
	}
	product := 1
	for i := 0; i < limit; i++ {
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
	fmt.Println("Part 1:", part1(points, lavaGraph, 1000))
	fmt.Println("Part 2:", part2(points, lavaGraph))
}
