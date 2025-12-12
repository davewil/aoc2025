package main

import (
	"fmt"
	"strconv"
	"strings"

	utils "github.com/davewil/aoc-utils"
)

type GridSize struct {
	Width, Height int
}

func parseInput(raw string) ([]string, error) {
	lines := utils.ReadLines(raw)

	total := 0
	for _, line := range lines {
		if strings.ContainsAny(line, "x") {

			parts := strings.Split(line, ":")
			sizeParts := strings.Split(parts[0], "x")
			width := 0
			height := 0
			fmt.Sscanf(sizeParts[0], "%d", &width)
			fmt.Sscanf(sizeParts[1], "%d", &height)
			gridSize := GridSize{Width: width, Height: height}
			fmt.Printf("Parsed grid size: %d\n", gridSize.Width*gridSize.Height)
			pieceParts := strings.Fields(parts[1])
			allPartsArea := 0
			for _, part := range pieceParts {
				p, _ := strconv.Atoi(part)
				allPartsArea += p * 7
			}
			fmt.Printf("All parts area: %d\n", allPartsArea)

			//assume that if all parts area is less than or equal to grid area, it might fit
			if allPartsArea <= gridSize.Width*gridSize.Height {
				total += 1
			}
		}
	}
	fmt.Printf("Total fitting grids: %d\n", total)
	// turns out that this is the answer
	return []string{}, nil
}
func part1(lines []string) int { return 0 }
func part2(lines []string) int { return 0 }

func main() {
	utils.LoadEnv()
	raw, err := utils.GetPuzzleInput(2025, 12)
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
