package utils

import (
	"container/heap"
)

// AStar finds the shortest path from start to target using the A* algorithm.
// T must be comparable.
// cost(a, b) returns the cost to move from a to b.
// heuristic(a) returns the estimated cost from a to the target.
// Returns the path, total cost, and true if found.
func AStar[T comparable](start T, isTarget func(T) bool, getNeighbours func(T) []T, cost func(T, T) int, heuristic func(T) int) ([]T, int, bool) {
	pq := make(PriorityQueue[T], 0)
	heap.Init(&pq)

	startItem := &Item[T]{
		Value:    start,
		Priority: 0 + heuristic(start),
		Cost:     0,
	}
	heap.Push(&pq, startItem)

	cameFrom := make(map[T]T)
	gScore := make(map[T]int)
	gScore[start] = 0

	visited := make(map[T]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item[T])
		u := current.Value

		if isTarget(u) {
			// Reconstruct path
			path := []T{}
			for curr := u; curr != start; curr = cameFrom[curr] {
				path = append([]T{curr}, path...)
			}
			path = append([]T{start}, path...)
			return path, gScore[u], true
		}

		visited[u] = true

		for _, v := range getNeighbours(u) {
			if visited[v] {
				continue
			}

			tentativeGScore := gScore[u] + cost(u, v)
			if existingG, ok := gScore[v]; !ok || tentativeGScore < existingG {
				cameFrom[v] = u
				gScore[v] = tentativeGScore
				fScore := tentativeGScore + heuristic(v)

				// We should update priority if it's already in PQ, but standard A* just adds it.
				// For simplicity and performance in many cases, we just push.
				// Duplicate entries are handled by the visited check or gScore check.
				heap.Push(&pq, &Item[T]{
					Value:    v,
					Priority: fScore,
					Cost:     tentativeGScore,
				})
			}
		}
	}

	return nil, 0, false
}

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T any] []*Item[T]

type Item[T any] struct {
	Value    T
	Priority int // The priority of the item in the queue.
	Cost     int // Actual cost so far (gScore)
	Index    int // The index of the item in the heap.
}

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the lowest priority (lowest fScore) so we use <
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
