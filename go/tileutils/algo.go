package tileutils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Point struct {
	x, y int
}

func FindNeighbors(t *Tile, startX, startY int, wg *sync.WaitGroup, exploreAdj chan [2]float64) {
	// fmt.Printf("Starting %d x: %d, y: %d\n", GetGoid(), startX, startY)
	defer wg.Done()
	// defer fmt.Printf("Finished %d\n", GetGoid())

	stack := make([]Point, 0, 1000)

	// startX and startY already tagged reachable by the main worker
	pushNeighbors(t, startX, startY, &stack)

	for len(stack) > 0 {
		// Pop
		idx := len(stack) - 1
		p := stack[idx]
		stack = stack[:idx]

		x, y := p.x, p.y

		t.Mutex.Lock()
		// Checking if already explored by another worker
		if t.Reachable[x][y] {
			t.Mutex.Unlock()
			continue
		}

		// Tagging as reachable
		t.Reachable[x][y] = true
		t.Mutex.Unlock()

		// If on edge of tile, start worker on other tile
		if x == 0 || x == MATRIX_SIZE-1 || y == 0 || y == MATRIX_SIZE-1 {
			coord := [2]float64{
				t.XLambertLL + t.CellSize*float64(x),
				t.YLambertLL + t.CellSize*float64(1000-y),
			}

			// Ajustements précis des coordonnées (comme dans ton code original)
			if x == 0 {
				coord[0] -= t.CellSize * 1.2
			}
			if x == MATRIX_SIZE-1 {
				coord[0] += t.CellSize * 1.2
			}
			if y == 0 {
				coord[1] += t.CellSize * 1.2
			}
			if y == MATRIX_SIZE-1 {
				coord[1] -= t.CellSize * 1.2
			}

			// Adding work to the main worker
			wg.Add(1)
			exploreAdj <- coord
		}

		pushNeighbors(t, x, y, &stack)
	}
}

func pushNeighbors(t *Tile, x, y int, stack *[]Point) {
	if x > 0 && t.PotentiallyReachable[x-1][y] {
		*stack = append(*stack, Point{x - 1, y})
	}
	if x < MATRIX_SIZE-1 && t.PotentiallyReachable[x+1][y] {
		*stack = append(*stack, Point{x + 1, y})
	}
	if y > 0 && t.PotentiallyReachable[x][y-1] {
		*stack = append(*stack, Point{x, y - 1})
	}
	if y < MATRIX_SIZE-1 && t.PotentiallyReachable[x][y+1] {
		*stack = append(*stack, Point{x, y + 1})
	}
}

func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}
