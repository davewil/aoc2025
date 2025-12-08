package main

import (
	_ "embed"
	"testing"
)

//go:embed input.txt
var benchmarkInput string

func BenchmarkPart1(b *testing.B) {
	raw := benchmarkInput
	for i := 0; i < b.N; i++ {
		points, graph, err := parseInput(raw)
		if err != nil {
			b.Fatalf("parseInput failed: %v", err)
		}
		result := part1(points, graph, 1000)
		if result == 0 {
			b.Fatal("part1 returned zero result during benchmark")
		}
	}
}
