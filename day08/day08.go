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

// JunctionData bundles the parsed points with the graph representation built from them.
type JunctionData struct {
	Points []Point3D
	Graph  *simple.WeightedUndirectedGraph
}

// pointNode wraps a Point3D so the coordinates travel with the Gonum node.
type pointNode struct {
	id    int64
	Point Point3D
}

func (n *pointNode) ID() int64 {
	return n.id
}

func newPointNode(id int64, p Point3D) *pointNode {
	return &pointNode{id: id, Point: p}
}

type edgeCandidate struct {
	u, v *pointNode
	w    float64
}

type unionFind struct {
	parent map[int64]int64
	rank   map[int64]int
}

func newUnionFind(nodes []*pointNode) *unionFind {
	uf := &unionFind{
		parent: make(map[int64]int64, len(nodes)),
		rank:   make(map[int64]int, len(nodes)),
	}
	for _, n := range nodes {
		uf.parent[n.ID()] = n.ID()
	}
	return uf
}

func (uf *unionFind) find(x int64) int64 {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *unionFind) union(a, b int64) bool {
	rootA := uf.find(a)
	rootB := uf.find(b)
	if rootA == rootB {
		return false
	}
	if uf.rank[rootA] < uf.rank[rootB] {
		rootA, rootB = rootB, rootA
	}
	uf.parent[rootB] = rootA
	if uf.rank[rootA] == uf.rank[rootB] {
		uf.rank[rootA]++
	}
	return true
}

func collectPointNodes(g *simple.WeightedUndirectedGraph) []*pointNode {
	var nodes []*pointNode
	if g == nil {
		return nodes
	}
	iter := g.Nodes()
	for iter.Next() {
		if n, ok := iter.Node().(*pointNode); ok {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func buildEdgeCandidates(nodes []*pointNode) []edgeCandidate {
	if len(nodes) < 2 {
		return nil
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
	return edges
}

func cloneJunctionGraph(src *simple.WeightedUndirectedGraph) *simple.WeightedUndirectedGraph {
	if src == nil {
		return nil
	}
	clone := simple.NewWeightedUndirectedGraph(0, math.Inf(1))
	nodeLookup := make(map[int64]*pointNode)
	nodes := src.Nodes()
	for nodes.Next() {
		if pn, ok := nodes.Node().(*pointNode); ok {
			copyNode := newPointNode(pn.ID(), pn.Point)
			nodeLookup[pn.ID()] = copyNode
			clone.AddNode(copyNode)
		}
	}
	edges := src.WeightedEdges()
	for edges != nil && edges.Next() {
		edge := edges.WeightedEdge()
		from := nodeLookup[edge.From().ID()]
		to := nodeLookup[edge.To().ID()]
		if from == nil || to == nil {
			continue
		}
		clone.SetWeightedEdge(simple.WeightedEdge{F: from, T: to, W: edge.Weight()})
	}
	return clone
}

func parsePoint3D(line string, lineNo int) (Point3D, error) {
	coords := strings.Split(line, ",")
	if len(coords) != 3 {
		return Point3D{}, fmt.Errorf("line %d: expected 3 values, got %d", lineNo, len(coords))
	}
	var values [3]int
	for i, coord := range coords {
		num, err := strconv.Atoi(strings.TrimSpace(coord))
		if err != nil {
			return Point3D{}, fmt.Errorf("line %d: invalid integer %q: %w", lineNo, coord, err)
		}
		values[i] = num
	}
	return Point3D{X: values[0], Y: values[1], Z: values[2]}, nil
}

func parseInput(raw string) (*JunctionData, error) {
	lines := utils.ReadLines(raw)
	data := &JunctionData{
		Points: make([]Point3D, 0, len(lines)),
		Graph:  simple.NewWeightedUndirectedGraph(0, math.Inf(1)),
	}
	var nextID int64
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		point, err := parsePoint3D(line, idx+1)
		if err != nil {
			return nil, err
		}
		data.Points = append(data.Points, point)
		node := newPointNode(nextID, point)
		nextID++
		data.Graph.AddNode(node)
	}
	return data, nil
}

func connectNearestPairs(g *simple.WeightedUndirectedGraph, count int) {
	if g == nil || count <= 0 {
		return
	}
	nodes := collectPointNodes(g)
	if len(nodes) < 2 {
		return
	}
	edges := buildEdgeCandidates(nodes)
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

func connectUntilSingleCircuit(g *simple.WeightedUndirectedGraph) (*edgeCandidate, bool) {
	if g == nil {
		return nil, false
	}
	nodes := collectPointNodes(g)
	if len(nodes) < 2 {
		return nil, false
	}
	edges := buildEdgeCandidates(nodes)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})
	uf := newUnionFind(nodes)
	components := len(nodes)
	var last edgeCandidate
	found := false
	for _, edge := range edges {
		if !uf.union(edge.u.ID(), edge.v.ID()) {
			continue
		}
		g.SetWeightedEdge(simple.WeightedEdge{F: edge.u, T: edge.v, W: edge.w})
		components--
		last = edge
		found = true
		if components == 1 {
			return &last, true
		}
	}
	return nil, found && components == 1
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
		sizes = append(sizes, visitComponent(n.ID(), g, visited))
	}
	return sizes
}

func visitComponent(start int64, g graph.Graph, visited map[int64]bool) int {
	size := 0
	stack := []int64{start}
	visited[start] = true
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
	return size
}

func part1(data *JunctionData, connectCount int) int {
	if data == nil || data.Graph == nil {
		return 0
	}
	workingGraph := cloneJunctionGraph(data.Graph)
	if workingGraph == nil {
		return 0
	}
	connectNearestPairs(workingGraph, connectCount)
	sizes := connectedComponentSizes(workingGraph)
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
func part2(data *JunctionData) int {
	if data == nil || data.Graph == nil {
		return 0
	}
	workingGraph := cloneJunctionGraph(data.Graph)
	if workingGraph == nil {
		return 0
	}
	lastEdge, ok := connectUntilSingleCircuit(workingGraph)
	if !ok || lastEdge == nil {
		return 0
	}
	return lastEdge.u.Point.X * lastEdge.v.Point.X
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 8)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	junctionData, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	start := time.Now()
	part1Result := part1(junctionData, 1000)
	fmt.Printf("Part 1: %d (took %s)\n", part1Result, time.Since(start))
	start = time.Now()
	part2Result := part2(junctionData)
	fmt.Printf("Part 2: %d (took %s)\n", part2Result, time.Since(start))
}
