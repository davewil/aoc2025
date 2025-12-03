package main

import (
	"testing"
)

var exampleInput = `987654321111111
811111111111119
234234234234278
818181911112111`

func TestPart1(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(lines)
	expected := 357
	if got != expected {
		t.Errorf("part1() = %d; expected %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}

	got := part2(lines)
	var expected int64 = 3121910778619
	if got != expected {
		t.Errorf("part2() = %d; expected %d", got, expected)
	}
}
