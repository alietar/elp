package algo

import "fmt"

type Matrix struct {
	Size int
	Data [][]float64
}

func NewMatrix(size int) Matrix {
	data := make([][]float64, size)

	for i := range data {
		data[i] = make([]float64, size)
	}

	return Matrix{
		Size: size,
		Data: data,
	}
}

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
