package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	utils "github.com/davewil/aoc-utils"
)

type Range struct {
	Start int
	End   int
}

func parseInput(input string) ([]Range, error) {
	cleaned := strings.Join(strings.Fields(input), "")
	data := strings.Split(cleaned, ",")
	ranges := make([]Range, 0, len(data))

	for _, r := range data {
		bounds := strings.Split(r, "-")
		if len(bounds) != 2 {
			return nil, fmt.Errorf("invalid range: %s", r)
		}
		start, err := strconv.Atoi(bounds[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(bounds[1])
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, Range{Start: start, End: end})
	}
	return ranges, nil

}

func part1(ranges []Range) int {
	total := 0
	for _, r := range ranges {
		for i := r.Start; i <= r.End; i++ {
			s := strconv.Itoa(i)
			if len(s)%2 != 0 {
				continue
			}
			mid := len(s) / 2
			if s[:mid] == s[mid:] {
				total += i
			}
		}
	}

	return total
}

func getDivisors(n int) []int {
	var result []int
	// We only go up to n/2 because a repeating chunk must be at most half the string
	for i := 1; i <= n/2; i++ {
		if n%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

func hasRepeatingPattern(s string, chunkLen int) bool {
	if chunkLen <= 0 || len(s)%chunkLen != 0 {
		return false
	}
	pattern := s[:chunkLen]
	for i := chunkLen; i < len(s); i += chunkLen {
		if s[i:i+chunkLen] != pattern {
			return false
		}
	}
	return true
}

func part2(ranges []Range) int {
	seen := make(map[int]bool)
	total := 0
	for _, r := range ranges {
		for i := r.Start; i <= r.End; i++ {
			if seen[i] {
				continue
			}

			s := strconv.Itoa(i)
			if len(s) < 2 {
				continue
			}

			divs := getDivisors(len(s))

			for _, d := range divs {
				if hasRepeatingPattern(s, d) {
					seen[i] = true
					total += i
					break
				}
			}
		}
	}
	return total
}

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 2)
	if err != nil {
		fmt.Println("Error fetching input:", err)
		return
	}
	data, err := parseInput(raw)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	start := time.Now()
	result1 := part1(data)
	duration1 := time.Since(start)
	fmt.Printf("Part 1: %d (took %v)\n", result1, duration1)

	start = time.Now()
	result2 := part2(data)
	duration2 := time.Since(start)
	fmt.Printf("Part 2: %d (took %v)\n", result2, duration2)
}
