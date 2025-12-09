package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	utils "github.com/davewil/aoc-utils"
)

type Point3D struct {
	X, Y, Z int
}

type edge struct {
	a, b int
	w    int
}

// JunctionSet keeps the parsed points alongside their sorted pairwise edges.
type JunctionSet struct {
	Points []Point3D
	Edges  []edge
}

type unionFind struct {
	parent []int
	rank   []int
}

func newUnionFind(size int) *unionFind {
	uf := &unionFind{
		parent: make([]int, size),
		rank:   make([]int, size),
	}
	for i := range uf.parent {
		uf.parent[i] = i
	}
	return uf
}

func (uf *unionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *unionFind) union(a, b int) bool {
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

func parseInput(raw string) (*JunctionSet, error) {
	lines := utils.ReadLines(raw)
	points := make([]Point3D, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		point, err := parsePoint3D(line, idx+1)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	edges := buildEdges(points)
	return &JunctionSet{Points: points, Edges: edges}, nil
}

func buildEdges(points []Point3D) []edge {
	n := len(points)
	if n < 2 {
		return nil
	}
	edges := make([]edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, edge{a: i, b: j, w: squaredDistance(points[i], points[j])})
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})
	return edges
}

func squaredDistance(a, b Point3D) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return dx*dx + dy*dy + dz*dz
}

func componentSizes(uf *unionFind) []int {
	counts := make(map[int]int)
	for i := range uf.parent {
		root := uf.find(i)
		counts[root]++
	}
	sizes := make([]int, 0, len(counts))
	for _, size := range counts {
		sizes = append(sizes, size)
	}
	return sizes
}

func part1(data *JunctionSet, connectCount int) int {
	if data == nil || len(data.Points) == 0 {
		return 0
	}
	if connectCount > len(data.Edges) {
		connectCount = len(data.Edges)
	}
	uf := newUnionFind(len(data.Points))
	for i := 0; i < connectCount; i++ {
		e := data.Edges[i]
		uf.union(e.a, e.b)
	}
	sizes := componentSizes(uf)
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

func part2(data *JunctionSet) int {
	if data == nil || len(data.Points) == 0 {
		return 0
	}
	uf := newUnionFind(len(data.Points))
	components := len(data.Points)
	for _, e := range data.Edges {
		if !uf.union(e.a, e.b) {
			continue
		}
		components--
		if components == 1 {
			return data.Points[e.a].X * data.Points[e.b].X
		}
	}
	return 0
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 8)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	junctions, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	start := time.Now()
	part1Result := part1(junctions, 1000)
	fmt.Printf("Part 1: %d (took %s)\n", part1Result, time.Since(start))
	start = time.Now()
	part2Result := part2(junctions)
	fmt.Printf("Part 2: %d (took %s)\n", part2Result, time.Since(start))
}
