package main

import (
	"testing"
)

var exampleInput = `aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out`

var exampleInput2 = `svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out`

func TestPart1(t *testing.T) {
	graph, err := parseInput(exampleInput)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part1(graph)
	expected := 5
	if got != expected {
		t.Errorf("part1() = %d; want %d", got, expected)
	}
}

func TestPart2(t *testing.T) {
	graph, err := parseInput(exampleInput2)
	if err != nil {
		t.Fatalf("parseInput error: %v", err)
	}
	got := part2(graph)
	expected := 2
	if got != expected {
		t.Errorf("part2() = %d; want %d", got, expected)
	}
}
