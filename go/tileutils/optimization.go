package tileutils

import "github.com/alietar/elp/go/gpsfiles"

func isSquareValid(position, used *[MATRIX_SIZE][MATRIX_SIZE]bool, startX, startY, size int) bool {
	for x := startX; x < startX+size; x++ {
		for y := startY; y < startY+size; y++ {
			if !(*position)[x][y] || (*used)[x][y] {
				return false
			}
		}
	}
	return true
}

func findMaxSquare(position, used *[MATRIX_SIZE][MATRIX_SIZE]bool, startX, startY int) int {
	size := 1

	for startX+size < MATRIX_SIZE && startY+size < MATRIX_SIZE {
		if !isSquareValid(position, used, startX, startY, size+1) {
			break
		}

		size++
	}

	return size
}

func markSquareUsed(used *[MATRIX_SIZE][MATRIX_SIZE]bool, startX, startY, size int) {
	for x := startX; x < startX+size; x++ {
		for y := startY; y < startY+size; y++ {
			(*used)[x][y] = true
		}
	}
}

func computeCenterLambert(x, y, size int, xLambertLL, yLambertLL, cellSize float64) (float64, float64) {
	cx := xLambertLL + (float64(x)+float64(size)/2.0)*cellSize
	cy := yLambertLL + (float64(MATRIX_SIZE)-(float64(y)+float64(size)/2.0))*cellSize

	return cx, cy
}

func addSquare(opt *[]Wgs84Square,
	x, y, size int,
	xLambertLL, yLambertLL, cellSize float64) {

	cx, cy := computeCenterLambert(x, y, size, xLambertLL, yLambertLL, cellSize)

	cx, cy, _ = gpsfiles.ConvertLambert93ToWgs84(cx, cy) // Conversion en coordonnées GPS WGS84

	*opt = append(*opt, Wgs84Square{
		Size:      size,
		CenterLng: cx,
		CenterLat: cy,
	})
}

func (t *Tile) ComputeOptimizedSquaresWgs() (opt []Wgs84Square) {
	var used [MATRIX_SIZE][MATRIX_SIZE]bool

	for x := range MATRIX_SIZE {
		for y := range MATRIX_SIZE {
			if t.Reachable[x][y] && !used[x][y] {
				size := findMaxSquare(t.Reachable, &used, x, y)

				markSquareUsed(&used, x, y, size)

				addSquare(&opt, x, y, size,
					t.XLambertLL, t.YLambertLL, 25)
			}
		}
	}

	return
}
