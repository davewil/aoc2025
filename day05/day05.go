package main

import (
    "fmt"
    utils "github.com/davewil/aoc-utils"
)

func parseInput(raw string) ([]string, error) { return []string{}, nil }
func part1(lines []string) int { return 0 }
func part2(lines []string) int { return 0 }

func main() {
    utils.LoadEnv()
    raw, err := utils.GetPuzzleInput(2025, 5)
    if err != nil { fmt.Println("Error fetching input:", err); return }
    lines, err := parseInput(raw)
    if err != nil { fmt.Println("Error parsing input:", err); return }
    fmt.Println("Part 1:", part1(lines))
    fmt.Println("Part 2:", part2(lines))
}
