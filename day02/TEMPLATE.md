# Day 02 Solution Template

This file documents the current skeleton in `day02.go`. It is intentionally minimal and mirrors the pattern established in Day 01, stopping short of any puzzle-specific logic until the official problem is released.

## Goals of the Template
- Provide a reproducible structure for every Advent of Code day.
- Make it trivial to drop in parsing + problem logic without refactoring.
- Keep early commits clean (pure scaffolding vs speculative implementation).

## Current Components (day02.go)
1. Imports: Only `fmt` and the shared `aoc-utils` package for input loading.
2. `parseInput(raw string) ([]string, error)`:
   - Stub: Returns an empty slice.
   - To implement: Convert the raw puzzle input into a slice of meaningful tokens/lines/records once the format is known.
3. `part1(lines []string) int` and `part2(lines []string) int`:
   - Stubs returning `0`.
   - To implement: Core logic for each part using the parsed input result.
4. `main()`:
   - Loads environment (for session cookie, etc.).
   - Fetches puzzle input for Day 02 using `GetPuzzleInput(2025, 2)`.
   - Calls `parseInput`, then `part1` and `part2`, printing their results.

## Rationale for Using `[]string`
Until the puzzle defines structure, the lowest-friction abstraction is a slice of strings. This avoids premature modeling (e.g., an `Instruction` or custom struct) and eliminates rewrite churn if the actual input format diverges.

## How to Extend When Puzzle Drops
1. Inspect raw input shape (lines, grids, CSV-like, groups separated by blanks, etc.).
2. Replace the body of `parseInput` with logic to populate a domain-specific type:
   - Option A: Keep `[]string` if line-based processing is sufficient.
   - Option B: Introduce a new struct type (e.g., `Record`, `Instruction`, `Entry`) and change return type accordingly.
3. Add focused unit tests:
   - Parsing: Assert correct element count + sample field extraction.
   - Part 1: Use provided example(s) to assert expected answer.
   - Part 2: Same pattern; add edge-case tests if logic branches significantly.
4. Optimize only after correctness (AoC inputs are small; clarity first).

## Suggested Implementation Checklist
- [ ] Replace `parseInput` stub.
- [ ] Add `exampleInput` constant to `day02_test.go` once sample is published.
- [ ] Implement `part1` and test against example.
- [ ] Implement `part2` and test against example.
- [ ] Run real input and record answers.
- [ ] (Optional) Add benchmarking if performance considerations arise.

## Testing Pattern (to mirror Day 01)
`day02_test.go` should look like:
```go
var exampleInput = `...` // Official sample once available

func TestPart1(t *testing.T) { /* parse -> part1 -> compare */ }
func TestPart2(t *testing.T) { /* parse -> part2 -> compare */ }
```
Keep expectations literal and tight; do not hide failures with helper layers early.

## Common Parse Patterns to Decide Between Later
- Line list: `strings.Split` on `\n`, trim blanks.
- Grouped blocks: split on double newlines, then split each block into lines.
- Grid: treat each line as a rune slice for coordinates.
- Token stream: regex to extract numeric + symbolic elements.

Pick the simplest thing that fits the actual data.

## Coding Style Notes
- Keep functions small and pure where possible (e.g., parse → transform → compute).
- Avoid global state; pass slices into parts explicitly.
- Preserve naming symmetry with other days (`part1`, `part2`, `parseInput`).

## Next Action
Await puzzle release; do not pre-build structure guesses. When ready, implement in-place without renaming existing function entry points.

---
This template is intentionally sparse; extend only with puzzle-driven needs.
