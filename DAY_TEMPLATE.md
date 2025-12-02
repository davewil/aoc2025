# Advent of Code Day Template

This document is a reusable template for any future Advent of Code day in this repository. Copy the structure and adapt only where the specific puzzle demands it.

## Directory Layout
```
/dayNN/
    dayNN.go
    dayNN_test.go
```
Where `NN` is a zero-padded day number (e.g. `03`, `14`).

## dayNN.go Skeleton
```go
package main

import (
    "fmt"
    utils "github.com/davewil/aoc-utils"
)

// parseInput parses the raw puzzle input into a slice of strings.
// Replace this with domain-specific parsing once the puzzle format is known.
func parseInput(raw string) ([]string, error) {
    return []string{}, nil
}

// part1 solves the first puzzle part. Implement once logic is known.
func part1(lines []string) int { return 0 }

// part2 solves the second puzzle part. Implement once logic is known.
func part2(lines []string) int { return 0 }

func main() {
    utils.LoadEnv()
    raw, err := utils.GetPuzzleInput(2025, NN) // Replace NN with the day number.
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
```

## dayNN_test.go Skeleton
```go
package main

import "testing"

// exampleInput is populated once the puzzle provides a sample.
var exampleInput = ``

func TestPart1(t *testing.T) {
    lines, err := parseInput(exampleInput)
    if err != nil {
        t.Fatalf("parseInput error: %v", err)
    }
    // Update expected after sample released.
    got := part1(lines)
    expected := 0
    if got != expected {
        t.Errorf("part1() = %d; want %d", got, expected)
    }
}

func TestPart2(t *testing.T) {
    lines, err := parseInput(exampleInput)
    if err != nil {
        t.Fatalf("parseInput error: %v", err)
    }
    got := part2(lines)
    expected := 0
    if got != expected {
        t.Errorf("part2() = %d; want %d", got, expected)
    }
}
```

## Usage Steps Per New Day
1. Create directory: `mkdir dayNN`.
2. Copy `dayNN.go` and `dayNN_test.go` from this template (replace NN).
3. Fetch input via existing `aoc-utils` (already used in Day 01).
4. Identify input format once puzzle unlocks.
5. Implement `parseInput` minimally (line split, regex, grouping, etc.).
6. Solve Part 1; add expected sample answer to test.
7. Solve Part 2; extend tests similarly.
8. Run: `go test -v ./dayNN/` and then `go run ./dayNN/` for real answers.

## Test Setup
Follow these steps to get a minimal test harness working quickly:

1. Create the test file: `touch dayNN/dayNN_test.go`.
2. Paste the skeleton from this template (keep `exampleInput` empty until the puzzle provides sample data).
3. Run initial compile check:
    ```bash
    go test -c ./dayNN/   # Produces test binary; ensures code compiles
    go test ./dayNN/      # Executes placeholder tests (should pass with zeros)
    ```
4. When the puzzle unlocks, update `exampleInput` with the provided sample block verbatim.
5. Replace each `expected := 0` with the sample’s published answers for Part 1 and Part 2 respectively.
6. If parsing requires transformation (numbers, grids, groups), implement `parseInput` first and re-run tests before writing solution logic.
7. Add edge-case tests only after baseline correctness (e.g. empty input, minimal valid input, large synthetic input if helpful).
8. Optional benchmarks (after solving):
    ```go
    func BenchmarkPart1(b *testing.B) {
         lines, _ := parseInput(exampleInput)
         for i := 0; i < b.N; i++ { _ = part1(lines) }
    }
    ```

### Test File Gotchas
- Avoid importing anything you do not use to keep compilation fast.
- Keep tests deterministic; do not fetch live input inside tests.
- Keep `exampleInput` small; large real inputs belong only in runtime execution (`main`).
- Ensure `parseInput` errors fail fast via `t.Fatalf` so broken parsing doesn’t silently produce wrong answers.

### Updating Expected Answers
Once you submit and receive the correct Part 1 answer, update its expected in the test immediately. Do the same for Part 2. This locks in regression protection if you refactor later.

### Running All Day Tests
From repository root you can run every solved day:
```bash
go test ./...
```
Use this before large refactors to ensure prior days remain green.

## Common Parse Patterns (Pick When Needed)
- Line list: `strings.Split(raw, "\n")` and trim blanks.
- Group blocks: `strings.Split(raw, "\n\n")` then per-block line split.
- Grid: convert each line to `[]rune` for coordinates.
- Token stream: use `regexp` to extract numbers or symbols.

## Principles
- Keep the first commit of each day strictly scaffolding.
- Avoid guessing structures until the puzzle format is visible.
- Keep naming consistent: `parseInput`, `part1`, `part2`.
- Prefer pure functions; no global mutable state.

## Optional Enhancements (After Correctness)
- Add benchmarks (`BenchmarkPart1`, `BenchmarkPart2`).
- Introduce small helper packages if repeated logic emerges (e.g. grid utilities).
- Record final answers in a `README.md` or results file.

## Checklist
- [ ] Directory created
- [ ] Skeleton files copied
- [ ] Sample input added to test
- [ ] parseInput implemented
- [ ] part1 passes sample
- [ ] part2 passes sample
- [ ] Real input answers printed

Copy this file forward as needed; avoid editing past days retroactively unless refactoring shared utilities.
