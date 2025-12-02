package main

import (
	"testing"

	utils "github.com/davewil/aoc-utils"
)

// To run these benchmarks:
// go test -bench=. -benchmem ./day02

var (
	benchInput string
	benchData  []Range
)

func init() {
	// Load input once for all benchmarks
	utils.LoadEnv()
	var err error
	benchInput, err = utils.GetPuzzleInput(2025, 2)
	if err != nil {
		panic(err)
	}
	benchData, err = parseInput(benchInput)
	if err != nil {
		panic(err)
	}
}

func BenchmarkPart1Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1(benchData)
	}
}

func BenchmarkPart1Concurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1Concurrent(benchData)
	}
}

func BenchmarkPart2Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2(benchData)
	}
}

func BenchmarkPart2Concurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2Concurrent(benchData)
	}
}
