package main

import (
	"os"
	"testing"
)

func BenchmarkPart1(b *testing.B) {
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Skip("input.txt missing, skipping benchmark")
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		data, err := parseInput(raw)
		if err != nil {
			b.Fatalf("parseInput failed: %v", err)
		}
		result := part1(data)
		if result == 0 {
			b.Fatal("part1 returned zero during benchmark")
		}
	}
}

func BenchmarkPart2(b *testing.B) {
	rawBytes, err := os.ReadFile("input.txt")
	if err != nil {
		b.Skip("input.txt missing, skipping benchmark")
	}
	raw := string(rawBytes)
	for i := 0; i < b.N; i++ {
		data, err := parseInput(raw)
		if err != nil {
			b.Fatalf("parseInput failed: %v", err)
		}
		result := part2(data)
		if result == 0 {
			b.Fatal("part2 returned zero during benchmark")
		}
	}
}
