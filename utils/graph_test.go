package utils

import (
	"testing"
)

func TestBFS(t *testing.T) {
	// Simple 3x3 grid
	// S . .
	// # # .
	// . . T
	start := Point{0, 0}
	target := Point{2, 2}
	walls := map[Point]bool{
		{0, 1}: true,
		{1, 1}: true,
	}

	isTarget := func(p Point) bool {
		return p == target
	}

	getNeighbours := func(p Point) []Point {
		neighbours := []Point{}
		dirs := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, d := range dirs {
			next := Point{p.X + d.X, p.Y + d.Y}
			if next.X >= 0 && next.X < 3 && next.Y >= 0 && next.Y < 3 && !walls[next] {
				neighbours = append(neighbours, next)
			}
		}
		return neighbours
	}

	path, found := BFS(start, isTarget, getNeighbours)
	if !found {
		t.Fatal("BFS path not found")
	}
	if len(path) != 5 {
		t.Errorf("BFS path length = %d; want 5. Path: %v", len(path), path)
	}
}

func TestAStar(t *testing.T) {
	// Same grid
	start := Point{0, 0}
	target := Point{2, 2}
	walls := map[Point]bool{
		{0, 1}: true,
		{1, 1}: true,
	}

	isTarget := func(p Point) bool {
		return p == target
	}

	getNeighbours := func(p Point) []Point {
		neighbours := []Point{}
		dirs := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, d := range dirs {
			next := Point{p.X + d.X, p.Y + d.Y}
			if next.X >= 0 && next.X < 3 && next.Y >= 0 && next.Y < 3 && !walls[next] {
				neighbours = append(neighbours, next)
			}
		}
		return neighbours
	}

	cost := func(a, b Point) int {
		return 1
	}

	heuristic := func(p Point) int {
		return ManhattanDistance(p, target)
	}

	path, totalCost, found := AStar(start, isTarget, getNeighbours, cost, heuristic)
	if !found {
		t.Fatal("AStar path not found")
	}
	if totalCost != 4 { // 4 steps
		t.Errorf("AStar totalCost = %d; want 4", totalCost)
	}
	expectedLen := 5
	if len(path) != expectedLen {
		t.Errorf("AStar path length = %d; want %d. Path: %v", len(path), expectedLen, path)
	}
}
