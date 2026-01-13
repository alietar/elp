package tileutils

const TILE_FOLDER_PATH = "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/"

const MATRIX_SIZE = 1000

type Tile struct {
	XLambertLL float64
	YLambertLL float64
	Lat        float64
	Long       float64

	Altitudes            *[MATRIX_SIZE][MATRIX_SIZE]float64
	PotentiallyReachable *[MATRIX_SIZE][MATRIX_SIZE]bool
	Reachable            *[MATRIX_SIZE][MATRIX_SIZE]bool
}

/*
func (m *Matrix) Show() {
	for i := 0; i < m.Size; i++ {
		for j := 0; j < m.Size; j++ {
			if m.Data[i][j] == 0 {
				fmt.Printf("       ")
			} else {
				fmt.Printf("%.2f ", m.Data[i][j])
			}
		}

		fmt.Printf("\n")
	}
}

func (m *Matrix) ShowPretty() {
	for i := 0; i < m.Size; i++ {
		for j := 0; j < m.Size; j++ {
			if m.Data[i][j] == 0 {
				fmt.Print("  ")
			} else {
				fmt.Print("██")
			}
		}

		fmt.Printf("\n")
	}
}

func (m *Matrix) ShowPrettyWithStart(x, y int) {
	for i := 0; i < m.Size; i++ {
		for j := 0; j < m.Size; j++ {
			if i == x && j == y {
				fmt.Print("\033[31m██\033[0m")
			} else {
				if m.Data[i][j] == 0 {
					fmt.Print("  ")
				} else {
					fmt.Print("██")
				}
			}
		}

		fmt.Printf("\n")
	}
}

func (m *Matrix) Resize(xCenter, yCenter, size int) (Matrix, int, int) {
	if size >= m.Size {
		return *m, xCenter, yCenter
	}

	clampedXCenter := Clamp(size/2, m.Size-size/2, xCenter)
	clampedYCenter := Clamp(size/2, m.Size-size/2, yCenter)

	newMatrix := NewMatrix(size)

	for i := range size {
		for j := range size {
			newMatrix.Data[i][j] = m.Data[(clampedXCenter-size/2)+i][(clampedYCenter-size/2)+j]
		}
	}

	return newMatrix, Min(size/2, xCenter), Min(size/2, yCenter)
}

func Clamp(min, max, a int) int {
	return Min(Max(a, min), max)
}

func Max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}*/
