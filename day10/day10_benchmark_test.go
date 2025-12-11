package main

import (
	"testing"

	utils "github.com/davewil/aoc-utils"
)

// To run these benchmarks:
// go test -bench=. -benchmem ./day10

var (
	benchData Input
)

func init() {
	// Load input once for all benchmarks
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 10)
	if err != nil {
		panic(err)
	}
	benchData, err = parseInput(raw)
	if err != nil {
		panic(err)
	}
}

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1(benchData)
	}
}

func BenchmarkPart2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2(benchData)
	}
}
