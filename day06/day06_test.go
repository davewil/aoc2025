package main

import (
	"testing"
)

var exampleInput = `123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  `

func TestPart1(t *testing.T) {
	lines, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(lines)
	expected := 4277556
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
	expected := 3263827
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
