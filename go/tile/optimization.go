package tile

type Square struct {
	X       int     // coin haut-gauche (colonne)
	Y       int     // coin haut-gauche (ligne)
	Size    int     // taille du carré
	CenterX float64 // Lambert X
	CenterY float64 // Lambert Y
}

type SquareOptimize struct {
	Squares []Square
}


func initUsed(n int) [][]bool {
	used := make([][]bool, n)
	for i := range used {
		used[i] = make([]bool, n)
	}
	return used
}

func isSquareValid(position, used [][]bool, startY, startX, size int) bool {
	for y := startY; y < startY+size; y++ {
		for x := startX; x < startX+size; x++ {
			if !position[y][x] || used[y][x] {
				return false
			}
		}
	}
	return true
}

func findMaxSquare(position, used [][]bool, startY, startX int) int {
	n := len(position)
	size := 1

	for startY+size < n && startX+size < n {
		if !isSquareValid(position, used, startY, startX, size+1) {
			break
		}
		size++
	}
	return size
}

func markSquareUsed(used [][]bool, startY, startX, size int) {
	for y := startY; y < startY+size; y++ {
		for x := startX; x < startX+size; x++ {
			used[y][x] = true
		}
	}
}

func computeCenterLambert(x, y, size int,
	xLambertLL, yLambertLL float64,
	cellSize int) (float64, float64) {

	cx := xLambertLL + (float64(x)+float64(size)/2.0)*float64(cellSize)
	cy := yLambertLL + (float64(y)+float64(size)/2.0)*float64(cellSize)
	return cx, cy
}

func addSquare(opt *SquareOptimize,
	x, y, size int,
	xLambertLL, yLambertLL float64,
	cellSize int) {

	cx, cy := computeCenterLambert(x, y, size, xLambertLL, yLambertLL, cellSize)

	opt.Squares = append(opt.Squares, Square{
		X:       x,
		Y:       y,
		Size:    size,
		CenterX: cx,
		CenterY: cy,
	})
}

func OptimizeSquares(position [][]bool,
	xLambertLL, yLambertLL float64,
	cellSize int) SquareOptimize {

	n := len(position)
	used := initUsed(n)

	var opt SquareOptimize

	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {

			if position[y][x] && !used[y][x] {

				size := findMaxSquare(position, used, y, x)

				markSquareUsed(used, y, x, size)

				addSquare(&opt, x, y, size,
					xLambertLL, yLambertLL, cellSize)
			}
		}
	}
	return opt
}



