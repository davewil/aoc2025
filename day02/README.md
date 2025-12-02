# Day 2: Red-Nosed Reports

## Running the Solution
To run the solution (which includes both sequential and concurrent versions with manual timing):

```bash
go run day02/day02.go day02/day02_chan.go
```

*Note: You must include both files because `day02.go` depends on the concurrent implementations in `day02_chan.go`.*

## Running Benchmarks
To run the scientific benchmarks comparing sequential vs. concurrent performance:

```bash
go test -bench=. -benchmem ./day02
```

## Files
*   `day02.go`: Main entry point and sequential implementation.
*   `day02_chan.go`: Concurrent implementations using worker pools.
*   `day02_benchmark_test.go`: Benchmark definitions.
