package main

import (
	"os"
	"testing"

	utils "github.com/davewil/aoc-utils"
)

func BenchmarkPart2Recursive(b *testing.B) {
	utils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		grid, _ := parseInput(raw)
		part2(grid)
	}
}

func BenchmarkPart3Iterative(b *testing.B) {
	utils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		grid, _ := parseInput(raw)
		part3(grid)
	}
}

func BenchmarkPart4Map(b *testing.B) {
	utils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		grid, _ := parseInput(raw)
		part4(grid)
	}
}

func BenchmarkPart5Optimized(b *testing.B) {
	utils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		grid, _ := parseInput(raw)
		part5(grid)
	}
}
