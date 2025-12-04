package main

import (
	"testing"
)

var exampleInput = `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

func TestPart1(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(lines)
	expected := 13
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
	expected := 0
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
