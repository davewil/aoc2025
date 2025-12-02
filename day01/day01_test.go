package main

import (
	"testing"
)

var exampleInput = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
`

func TestPart1(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(lines)
	expected := 3
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
	expected := 6
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
