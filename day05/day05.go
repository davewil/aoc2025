package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	utils "github.com/davewil/aoc-utils"
)

type Range struct {
	Start, End int
}

type Input struct {
	Ranges  []Range
	Numbers []int
}

func (r Range) Len() int {
	return r.End - r.Start + 1
}

func parseInput(raw string) (Input, error) {
	lines := utils.ReadLines(raw)
	ranges := []Range{}
	numbers := []int{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.Contains(line, "-") {
			rangeParts := strings.Split(line, "-")
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])
			ranges = append(ranges, Range{Start: start, End: end})
		} else {
			num, _ := strconv.Atoi(line)
			numbers = append(numbers, num)
		}
	}
	ranges = mergeRanges(ranges)
	return Input{Ranges: ranges, Numbers: numbers}, nil
}

func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return ranges
	}

	// Sort by Start
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	merged := []Range{ranges[0]}
	for _, current := range ranges[1:] {
		last := &merged[len(merged)-1]
		// Check for overlap
		if current.Start <= last.End {
			if current.End > last.End {
				last.End = current.End
			}
		} else {
			merged = append(merged, current)
		}
	}
	return merged
}

func part1(input Input) int {
	count := 0
	for _, num := range input.Numbers {
		for _, rng := range input.Ranges {
			if num >= rng.Start && num <= rng.End {
				count++
				break
			}
		}
	}
	return count
}
func part2(input Input) int {
	count := 0
	for _, rng := range input.Ranges {
		count += rng.Len()
	}
	return count
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 5)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	lines, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
}
