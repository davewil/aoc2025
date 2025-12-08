package main

import (
	"testing"
)

var exampleInput = `162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689`

func TestPart1(t *testing.T) {
	points, g, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(points, g, 10)
	expected := 40
	if got != expected {
		t.Errorf("part1() = %d; want %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	points, g, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part2(points, g)
	expected := 0
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
