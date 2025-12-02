package main

import (
	"testing"
)

var exampleInput = `11-22,95-115,998-1012,1188511880-1188511890,222220-222224,
1698522-1698528,446443-446449,38593856-38593862,565653-565659,
824824821-824824827,2121212118-2121212124`

var exampleInputPart2 = `998-1012`

func TestPart1(t *testing.T) {
	ranges, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(ranges)
	expected := 1227775554
	if got != expected {
		t.Errorf("part1() = %d; want %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	ranges, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part2(ranges)
	expected := 4174379265
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
