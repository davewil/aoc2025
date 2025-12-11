package main

import (
	"testing"
	"aoc2025/day10/parser"
)

var exampleInput = `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}`

func TestPart1(t *testing.T) {
	lines, err := parser.ParseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(lines)
	expected := 7
	if got != expected {
		t.Errorf("part1() = %d; want %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	lines, err := parser.ParseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part2(lines)
	expected := 33
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
