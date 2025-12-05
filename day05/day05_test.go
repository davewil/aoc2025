package main

import (
	"reflect"
	"testing"
)

var exampleInput = `3-5
10-14
16-20
12-18

1
5
8
11
17
32`

func TestPart1(t *testing.T) {
	input, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	ranges := input.Ranges
	numbers := input.Numbers

	expectedRanges := []Range{
		{Start: 3, End: 5},
		{Start: 10, End: 20},
	}
	if !reflect.DeepEqual(ranges, expectedRanges) {
		t.Errorf("parseInput() = %v; want %v", ranges, expectedRanges)
	}
	expectedNumbers := []int{1, 5, 8, 11, 17, 32}
	if !reflect.DeepEqual(numbers, expectedNumbers) {
		t.Errorf("parseInput() = %v; want %v", numbers, expectedNumbers)
	}
	got := part1(Input{Ranges: ranges, Numbers: numbers})

	expected := 3
	if got != expected {
		t.Errorf("part1() = %d; want %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	input, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part2(input)
	expected := 14
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
