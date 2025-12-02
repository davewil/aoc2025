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
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range: %s", r)
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(parts[1])
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
			mid := len(s) / 2
			left := s[:mid]
			right := s[mid:]
			if left == right {
				total += i
			}
		}
	}

	return total
}

var factorsCache = make(map[int][]int)

func factors(n int) []int {
	if facs, ok := factorsCache[n]; ok {
		return facs
	}

	var result []int
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			result = append(result, i)
			if i != n/i && n/i != n {
				result = append(result, n/i)
			}
		}

	}
	factorsCache[n] = result
	return result
}

func chunkString(s string, size int) []string {
	if size <= 0 {
		return nil
	}

	var chunks []string
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func allSame(slice []string) bool {
	if len(slice) == 0 {
		return true
	}
	first := slice[0]
	for _, s := range slice[1:] {
		if s != first {
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
			s := strconv.Itoa(i)
			if len(s) == 1 {
				continue
			}
			facs := factors(len(s))

			for _, f := range facs {
				chunks := chunkString(s, f)
				if allSame(chunks) && !seen[i] {
					seen[i] = true
					total += i
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
