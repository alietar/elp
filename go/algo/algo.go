package algo

import (
	"fmt"
	"sync"
)

func (m *Matrix) FindNeighbors(startX, startY int, wg *sync.WaitGroup, results chan [][]float64, exploreAdj chan [2]float64) {
	// Initialization of the matrices result and visited
	result := make([][]float64, m.Size)
	visited := make([][]bool, m.Size)

	for i := range m.Size {
		result[i] = make([]float64, m.Size)
		visited[i] = make([]bool, m.Size)
	}

	/*return Matrix{
		Size: m.Size,
		Data: m.findNeighborsRecursive(result, visited, startX, startY),
	}*/

	results <- m.findNeighborsRecursive(result, visited, startX, startY, exploreAdj, wg)

	wg.Done()
}

func (m *Matrix) findNeighborsRecursive(
	result [][]float64,
	visited [][]bool,
	x, y int,
	exploreAdj chan [2]float64,
	wg *sync.WaitGroup,
) [][]float64 {
	if m.Data[x][y] == 0 {
		return result
	}

	visited[x][y] = true
	result[x][y] = m.Data[x][y]

	if x == 0 || x == m.Size-1 || y == 0 || y == m.Size-1 {
		var coord [2]float64
		coord[0] = m.LambertX + float64((x-m.StartX)*25)
		coord[1] = m.LambertY + float64((y-m.StartY)*25)

		if x == 0 {
			coord[0] -= 25
		}
		if x == m.Size-1 {
			coord[0] += 25
		}
		if y == 0 {
			coord[1] -= 25
		}
		if y == m.Size-1 {
			coord[1] += 25
		}

		fmt.Println("Reached bounderies")
		fmt.Printf("x: %d, oldLambert: %f, newLambert: %f\n", x, m.LambertX, coord[0])
		fmt.Printf("y: %d, oldLambert: %f, newLambert: %f\n", y, m.LambertY, coord[1])

		wg.Add(1)
		exploreAdj <- coord

		return result
	}

	if !visited[x-1][y] && m.Data[x-1][y] != 0 {
		result = m.findNeighborsRecursive(result, visited, x-1, y, exploreAdj, wg)
	}

	if !visited[x+1][y] && m.Data[x+1][y] != 0 {
		result = m.findNeighborsRecursive(result, visited, x+1, y, exploreAdj, wg)
	}

	if !visited[x][y-1] && m.Data[x][y-1] != 0 {
		result = m.findNeighborsRecursive(result, visited, x, y-1, exploreAdj, wg)
	}
	if !visited[x][y+1] && m.Data[x][y+1] != 0 {
		result = m.findNeighborsRecursive(result, visited, x, y+1, exploreAdj, wg)
	}

	return result
}
