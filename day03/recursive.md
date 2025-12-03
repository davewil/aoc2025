# Recursive Approach for Maximum Joltage

## Overview

An alternative recursive algorithm for finding the maximum number by selecting a fixed number of digits from an input sequence.

## Algorithm

Instead of removing digits, this approach **selects** the best digit at each position:

1. Find the maximum digit in a valid range (ensuring enough digits remain)
2. Use that digit in the current position
3. Recursively solve for the remaining positions with the digits after the selected one

## Implementation

```go
import (
	"math"
	"slices"
)

func getHighestJoltage(bank []int) int64 {
	return maxJoltage(bank, 12)
}

func maxJoltage(bank []int, digits int) int64 {
	if digits == 1 {
		return int64(slices.Max(bank))
	}

	// Find max digit in the range where we can still get enough remaining digits
	maxDigit := slices.Max(bank[:len(bank)-(digits-1)])
	idx := slices.Index(bank, maxDigit)

	return int64(maxDigit)*int64(math.Pow10(digits-1)) + maxJoltage(bank[idx+1:], digits-1)
}
```

That's it! No helper functions needed - `slices.Max()`, `slices.Index()`, and `math.Pow10()` handle everything.

## Example

For input `818181` selecting 3 digits:

1. **First digit:** Find max in `8181` (need 2 more after) → `8` at index 0
   - Result so far: `8 * 100`
   - Recurse on `[1,8,1,8,1]` with 2 digits

2. **Second digit:** Find max in `1818` (need 1 more after) → `8` at index 1
   - Result so far: `8 * 10`
   - Recurse on `[1,8,1]` with 1 digit

3. **Third digit:** Max of `[1,8,1]` → `8`
   - Result: `8`

4. **Final:** `800 + 80 + 8 = 888`

## Complexity

- **Time:** O(n × d) where n is input length, d is number of digits to select
- **Space:** O(d) for recursion stack + O(n) for slice copies on each recursion level

## Comparison with Iterative Approach

### Advantages
- Much more concise and elegant
- Easier to reason about the algorithm
- Directly expresses "pick the best at each position"
- Uses standard library functions

### Disadvantages
- Creates slice copies on each recursive call (expensive in Go)
- Higher memory usage
- Go doesn't optimize tail recursion
- Less idiomatic for Go (Go culture prefers iteration)
- More GC pressure from allocations

## When to Use

This approach is clearer when:
- You want to understand the algorithm conceptually
- Performance isn't critical
- You're prototyping or teaching the concept

The iterative approach is preferred when:
- Performance matters
- Writing production Go code
- Following Go idioms and best practices
- Minimizing allocations
