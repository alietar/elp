package tileutils

import (
	"fmt"
	"sync"
)

func FindNeighbors(t *Tile, startX, startY int, wg *sync.WaitGroup, results chan *Tile, exploreAdj chan [2]float64) {
	results <- t

	findNeighborsRecursive(t, startX, startY, exploreAdj, wg)
	wg.Done()
}

func findNeighborsRecursive(
	t *Tile,
	x, y int,
	exploreAdj chan [2]float64,
	wg *sync.WaitGroup,
) {
	if !t.PotentiallyReachable[x][y] {
		return
	}
	if t.Reachable[x][y] { // mesure de sécurité ??
		return
	}

	t.Mutex.Lock()
	if t.Reachable[x][y] { // Double check après le lock
		fmt.Println("Quitting")
		t.Mutex.Unlock()
		return
	}
	t.Reachable[x][y] = true
	t.Mutex.Unlock()

	if x == 0 || x == MATRIX_SIZE-1 || y == 0 || y == MATRIX_SIZE-1 {
		coord := [2]float64{
			t.XLambertLL + t.CellSize*float64(x),
			t.YLambertLL + t.CellSize*float64(1000-y),
		}

		if x == 0 {
			coord[0] -= t.CellSize + t.CellSize/5
		}
		if x == MATRIX_SIZE-1 {
			coord[0] += t.CellSize + t.CellSize/5
		}
		if y == 0 {
			coord[1] += t.CellSize + t.CellSize/5
		}
		if y == MATRIX_SIZE-1 {
			coord[1] -= t.CellSize + t.CellSize/5
		}

		// fmt.Println("Reached bounderies")
		// fmt.Printf("x: %d, newLambert: %f\n", x, coord[0])
		// fmt.Printf("y: %d, newLambert: %f\n", y, coord[1])

		wg.Add(1)
		exploreAdj <- coord
	}

	if x > 0 {
		findNeighborsRecursive(t, x-1, y, exploreAdj, wg)
	}

	if x < MATRIX_SIZE-1 {
		findNeighborsRecursive(t, x+1, y, exploreAdj, wg)
	}

	if y > 0 {
		findNeighborsRecursive(t, x, y-1, exploreAdj, wg)
	}
	if y < MATRIX_SIZE-1 {
		findNeighborsRecursive(t, x, y+1, exploreAdj, wg)
	}
}
