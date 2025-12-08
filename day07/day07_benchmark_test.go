package main

import (
	"os"
	"testing"

	aocutils "github.com/davewil/aoc-utils"
)

func BenchmarkPart1(b *testing.B) {
	aocutils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		input, _ := parseInput(raw)
		part1(input)
	}
}

func BenchmarkPart2(b *testing.B) {
	aocutils.LoadEnv()
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		input, _ := parseInput(raw)
		part2(input)
	}
}
