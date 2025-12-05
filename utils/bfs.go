package utils

// BFS finds the shortest path from start to a target using Breadth-First Search.
// T must be comparable (usable as a map key).
// It requires:
// - start: the starting node.
// - isTarget: a function that returns true if a node is the target.
// - getNeighbours: a function that returns the neighbors of a node.
// Returns the path (including start and target) and true if found, or nil and false if not.
func BFS[T comparable](start T, isTarget func(T) bool, getNeighbours func(T) []T) ([]T, bool) {
	queue := []T{start}
	visited := make(map[T]bool)
	visited[start] = true
	parent := make(map[T]T)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if isTarget(current) {
			// Reconstruct path
			path := []T{}
			for curr := current; curr != start; curr = parent[curr] {
				path = append([]T{curr}, path...)
			}
			path = append([]T{start}, path...)
			return path, true
		}

		for _, neighbour := range getNeighbours(current) {
			if !visited[neighbour] {
				visited[neighbour] = true
				parent[neighbour] = current
				queue = append(queue, neighbour)
			}
		}
	}
	return nil, false
}
