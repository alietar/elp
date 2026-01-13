package tileutils

import (
	"fmt"
	"sync"
)

func FindNeighbors(t *Tile, startX, startY int, wg *sync.WaitGroup, results chan *Tile, exploreAdj chan [2]float64) {
	findNeighborsRecursive(t, startX, startY, exploreAdj, wg)

	results <- t

	wg.Done()
}

func findNeighborsRecursive(
	t *Tile,
	x, y int,
	exploreAdj chan [2]float64,
	wg *sync.WaitGroup,
) {
	if t.PotentiallyReachable[x][y] == false {
		return
	}

	t.Reachable[x][y] = true

	if x == 0 || x == MATRIX_SIZE-1 || y == 0 || y == MATRIX_SIZE-1 {
		coord := [2]float64{
			t.XLambertLL + float64(x*25),
			t.YLambertLL + float64((1000-y)*25),
		}

		if x == 0 {
			coord[0] -= 30
		}
		if x == MATRIX_SIZE-1 {
			coord[0] += 30
		}
		if y == 0 {
			coord[1] -= 30
		}
		if y == MATRIX_SIZE-1 {
			coord[1] += 30
		}

		fmt.Println("Reached bounderies")
		fmt.Printf("x: %d, newLambert: %f\n", x, coord[0])
		fmt.Printf("y: %d, newLambert: %f\n", y, coord[1])

		wg.Add(1)
		exploreAdj <- coord
	}

	if x > 0 && !t.Reachable[x-1][y] && t.PotentiallyReachable[x-1][y] {
		findNeighborsRecursive(t, x-1, y, exploreAdj, wg)
	}

	if x < MATRIX_SIZE-1 && !t.Reachable[x+1][y] && t.PotentiallyReachable[x+1][y] {
		findNeighborsRecursive(t, x+1, y, exploreAdj, wg)
	}

	if y > 0 && !t.Reachable[x][y-1] && t.PotentiallyReachable[x][y-1] {
		findNeighborsRecursive(t, x, y-1, exploreAdj, wg)
	}
	if y < MATRIX_SIZE-1 && !t.Reachable[x][y+1] && t.PotentiallyReachable[x][y+1] {
		findNeighborsRecursive(t, x, y+1, exploreAdj, wg)
	}
}
