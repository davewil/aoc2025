# Day 3: Recursive vs Iterative Joltage

This day's solution explores two different algorithmic approaches to finding the highest joltage number by selecting 12 digits from a sequence.

## Running the Solution

```bash
go run day03/day03.go
```

Output:
*   **Part 1**: Highest pair logic.
*   **Part 2**: Iterative "Greedy Removal" approach.
*   **Part 3**: Recursive "Greedy Selection" approach.

## Running Benchmarks

To compare the performance of the iterative vs. recursive approaches:

```bash
go test -bench=. -benchmem ./day03
```

### Benchmark Results

The recursive solution (Part 3) proved to be significantly faster and more memory-efficient because it utilizes Go slices (zero-copy views) rather than copying the array for destructive modification.

```
BenchmarkPart1-8             930           1127782 ns/op          499171 B/op       6764 allocs/op
BenchmarkPart2-8            1938            920130 ns/op          179200 B/op        200 allocs/op
BenchmarkPart3-8           10000            127532 ns/op               0 B/op          0 allocs/op
```

*   **Part 2 (Iterative)**: ~0.92ms/op, 200 allocs/op
*   **Part 3 (Recursive)**: ~0.13ms/op, 0 allocs/op

## Files

*   `day03.go`: Main implementation containing all 3 parts.
*   `day03_test.go`: Unit tests for correctness.
*   `day03_benchmark_test.go`: Performance benchmarks.
