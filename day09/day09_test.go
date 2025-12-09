package main

import (
	"testing"
)

var exampleInput = `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

func TestPart1(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(lines)
	expected := 50
	if got != expected {
		t.Errorf("part1() = %d; want %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part2(lines)
	expected := 24
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
