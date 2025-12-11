# Advent of Code 2025 (Go)

Solutions for Advent of Code 2025 written in Go.

## Prerequisites

- Go 1.25+
- Make (optional, for using the Makefile)

## Usage

### Running a Solution

To run the solution for a specific day (e.g., Day 10):

```bash
# Using Make
make run day=10

# Using Go directly
go run ./day10
```

### Testing

To run all unit tests:

```bash
# Using Make
make test

# Using Go directly
go test ./...
```

### Benchmarking

To run benchmarks:

```bash
# Using Make
make bench

# Using Go directly
go test -bench=. -benchmem ./...
```

## Project Structure

- `dayXX/`: Contains the solution and tests for each day.
- `utils/`: Shared utility functions.
- `Makefile`: Helper commands for testing and running.
