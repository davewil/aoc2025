# What is Happening in Day 10?

This document explains the approach used to solve Part 2 of Day 10, which involves finding the minimum number of button presses to achieve specific target values ("joltages") on a set of lights.

## The Problem: Integer Linear Programming

The puzzle asks us to find a combination of button presses.
*   We have $N$ buttons.
*   We have $M$ lights.
*   Each button, when pressed, adds 1 to the value of specific lights.
*   Each light must end up at a specific integer value (target joltage).
*   We want to **minimize the total number of button presses**.

Mathematically, this is a system of linear equations:

$$
\begin{cases}
b_0 \cdot A_{0,0} + b_1 \cdot A_{1,0} + \dots = \text{Target}_0 \\
b_0 \cdot A_{0,1} + b_1 \cdot A_{1,1} + \dots = \text{Target}_1 \\
\dots
\end{cases}
$$

Where:
*   $b_i$ is the number of times button $i$ is pressed (must be an integer $\ge 0$).
*   $A_{i,j}$ is 1 if button $i$ affects light $j$, and 0 otherwise.

Since we need integer solutions and want to minimize a linear objective function ($\sum b_i$), this is an **Integer Linear Programming (ILP)** problem.

## The Tool: Z3 SMT Solver

To solve this efficiently, we use **Z3**, a high-performance theorem prover from Microsoft Research. Z3 is an **SMT (Satisfiability Modulo Theories)** solver.

*   **SAT Solvers** figure out if a boolean formula (e.g., `(A or B) and (not A)`) can be true.
*   **SMT Solvers** extend this to "theories" like arithmetic, bitvectors, and arrays. They can answer questions like: "Is there a set of integers $x, y$ such that $x > 0, y > 0, 2x + y = 10$?"

Z3 is incredibly powerful at pruning the search space for these kinds of combinatorial problems, far faster than a brute-force BFS or simple backtracking approach.

## The Solution Strategy

### 1. Modeling
For each row of the puzzle, we create a Z3 `Solver` instance and define:
*   **Variables**: One integer variable for each button (`b0`, `b1`, etc.).
*   **Constraints**:
    1.  `b_i >= 0` for all buttons (cannot press a button negative times).
    2.  For each light, the sum of the variables for buttons affecting it must equal the target joltage.

### 2. Optimization (Model-Guided Search)
Z3 has an `Optimize` module, but for this specific problem, a custom binary search approach proved faster and more stable.

We want to find the minimum `TotalPresses` ($\sum b_i$). We know the answer is between 0 and the sum of all target joltages (upper bound).

We perform a binary search on the answer:
1.  Pick a candidate limit `K`.
2.  Add a temporary constraint: `TotalPresses <= K`.
3.  Ask Z3: `Check()`. Is this possible?
    *   **If SAT (Satisfiable):** Z3 found a valid combination with $\le K$ presses. We store this as the best answer so far and try to find an even smaller one (lower the upper bound).
    *   **If UNSAT (Unsatisfiable):** It's impossible to solve the puzzle with $\le K$ presses. We must increase the lower bound.

**The "Model-Guided" Twist:**
Instead of a standard binary search (splitting the range in half), when Z3 finds a solution, we ask it for the **actual value** of `TotalPresses` in that solution.
*   If we asked for $\le 50$ and Z3 found a solution with $30$ presses, we don't just set the new upper bound to $49$. We set it to $29$. This allows us to skip huge chunks of the search space.

### 3. Performance Optimizations

To make this run in milliseconds rather than seconds, several engineering tricks were applied:

*   **Incremental Solving (`Push`/`Pop`)**:
    We don't recreate the solver for every step of the binary search. We use `solver.Push()` to create a backtracking point, add the `TotalPresses <= K` constraint, check it, and then `solver.Pop()` to remove it. This keeps all the "hard" constraints (the linear equations) in memory and allows Z3 to reuse learned lemmas.

*   **Parallelism**:
    The puzzle input has many independent rows. We spin up a worker pool (one per CPU core) to solve multiple rows simultaneously.

*   **Variable Pooling**:
    Creating Z3 variables involves CGO (C-Go) calls and string formatting, which is slow. We pre-allocate a pool of Z3 integer variables for each worker and reuse them across rows, resetting the solver state in between.

*   **Variadic Addition**:
    Instead of building a sum tree like `((a + b) + c) + d`, we pass all terms to Z3's addition function at once (`Sum(a, b, c, d)`). This creates a flatter internal representation that is easier for the solver to manipulate.

## Summary

By translating the puzzle into a system of constraints and using a state-of-the-art mathematical solver (Z3), we transform a complex search problem into a series of efficient feasibility checks. The combination of algorithmic improvements (model-guided search) and low-level optimizations (pooling, parallelism) results in a solution that is orders of magnitude faster than naive approaches.
