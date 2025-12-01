package main

import (
	"testing"
)

var exampleInput = ``

func TestPart1(t *testing.T) {
	input, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput() error = %v", err)
	}

	result := part1(input)
	expected := 0

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
	expected := 0

	if result != expected {
		t.Errorf("part2() = %d; want %d", result, expected)
	}
}
