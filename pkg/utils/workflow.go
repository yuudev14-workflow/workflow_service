package utils

func IsAcyclicGraph(graph map[string][]string) bool {
	visit := make(Set[string])
	stack := make(Set[string])

	var dfs func(node string) bool
	dfs = func(node string) bool {
		if stack.Has(node) {
			return true
		}
		if visit.Has(node) {
			return false
		}

		visit.Add(node)
		stack.Add(node)

		for _, neighbor := range graph[node] {
			if dfs(neighbor) {
				return true
			}
		}

		stack.Remove(node)

		return false
	}

	for node := range graph {
		if dfs(node) {
			return true
		}
	}

	return false

}
