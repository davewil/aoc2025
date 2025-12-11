package main

import (
	"os"
	"testing"

	utils "github.com/davewil/aoc-utils"
)

func BenchmarkPart1(b *testing.B) {
	utils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		graph, _ := parseInput(raw)
		part1(graph)
	}
}

func BenchmarkPart2(b *testing.B) {
	utils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		graph, _ := parseInput(raw)
		part2(graph)
	}
}
