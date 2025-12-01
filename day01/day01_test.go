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
	input, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput() error = %v", err)
	}

	result := part1(input)
	expected := 3

	if result != expected {
		t.Errorf("part1() = %d; want %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	input, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput() error = %v", err)
	}

	result := part2(input)
	expected := 6

	if result != expected {
		t.Errorf("part2() = %d; want %d", result, expected)
	}
}
