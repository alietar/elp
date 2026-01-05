package algo

func ExploreAdjacentTile(startI, startJ) {
	// Get the adjacent matrix
	// Run Find neighbors on it
}

func (m *Matrix) FindNeighbors(startX, startY int, done chan bool, exploreAdjacentTile chan [2]int) Matrix {
	// Initialization of the matrices result and visited
	result := make([][]float64, m.Size)
	visited := make([][]bool, m.Size)

	for i := range m.Size {
		result[i] = make([]float64, m.Size)
		visited[i] = make([]bool, m.Size)
	}

	return Matrix{
		Size: m.Size,
		Data: m.findNeighborsRecursive(result, visited, startX, startY),
	}
}

func (m *Matrix) findNeighborsRecursive(
	result [][]float64,
	visited [][]bool,
	startX, startY int,
) [][]float64 {
	x := startX
	y := startY

	if m.Data[x][y] == 0 {
		return result
	}

	visited[x][y] = true
	result[x][y] = m.Data[x][y]

	if x == 0 {

	}

	if x > 0 && !visited[x-1][y] && m.Data[x-1][y] != 0 {
		result = m.findNeighborsRecursive(result, visited, x-1, y)
	}

	if x < m.Size-1 && !visited[x+1][y] && m.Data[x+1][y] != 0 {
		result = m.findNeighborsRecursive(result, visited, x+1, y)
	}

	if y > 0 && !visited[x][y-1] && m.Data[x][y-1] != 0 {
		result = m.findNeighborsRecursive(result, visited, x, y-1)
	}
	if y < m.Size-1 && !visited[x][y+1] && m.Data[x][y+1] != 0 {
		result = m.findNeighborsRecursive(result, visited, x, y+1)
	}

	return result
}
