# 🎄 Advent of Code 2025 🎄

![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)
![CI Status](https://github.com/davewil/aoc2025/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/davewil/aoc2025)](https://goreportcard.com/report/github.com/davewil/aoc2025)
![Completion](https://img.shields.io/badge/stars-20%20%2F%2050-yellow)

My solutions for [Advent of Code 2025](https://adventofcode.com/2025), written in **Go**.

This repository focuses on **clean code**, **performance**, and exploring advanced algorithms (like Z3 SMT solving).

## 🚀 Highlights

| Day | Title | Approach / Key Tech |
|:---:|:---|:---|
| **10** | **The Z3 Solver** | Used **Microsoft Z3** (SMT Solver) via CGO to solve an Integer Linear Programming problem. Includes parallel execution and model-guided binary search. [Read the Deep Dive](./day10/WHAT_HAPPENING.md) |
| **09** | **Tile Geometry** | Geometric algorithms to find maximum areas and optimal rectangles from coordinate sets. |
| **03** | **Recursion** | Explored recursive patterns and memoization. |

## 📊 Progress

| Day | Part 1 | Part 2 | Time (approx) | Tags |
|:---:|:---:|:---:|:---:|:---|
| 01 | ⭐ | ⭐ | < 1ms | `parsing` `math` |
| 02 | ⭐ | ⭐ | < 1ms | `strings` |
| 03 | ⭐ | ⭐ | 15ms | `recursion` |
| 04 | ⭐ | ⭐ | 5ms | `grid` |
| 05 | ⭐ | ⭐ | 2ms | `simulation` |
| 06 | ⭐ | ⭐ | 10ms | `optimization` |
| 07 | ⭐ | ⭐ | 8ms | `graph` |
| 08 | ⭐ | ⭐ | 12ms | `parsing` |
| 09 | ⭐ | ⭐ | 45ms | `simulation` |
| 10 | ⭐ | ⭐ | 250ms | `z3` `smt` `cgo` `parallel` |
| 11 |   |   | | |
| 12 |   |   | | |

## 🛠️ Usage

This project uses a `Makefile` for convenience.

### Run a Solution
```bash
make run day=10
```

### Run Tests
```bash
make test
```

### Run Benchmarks
```bash
make bench
```

## 📂 Structure

*   `dayXX/`: Individual solutions.
    *   `parser/`: (Optional) Complex parsing logic separated into its own package.
*   `utils/`: Shared helpers (Grid, Point, Math).
*   `Makefile`: Task runner.

---
*Created by [davewil](https://github.com/davewil)*
