package algo

import (
	"fmt"
)

const TAILLE_MATRICE int = 5

func afficherMatrice(matrice [TAILLE_MATRICE][TAILLE_MATRICE]float64) {
	for i := 0; i < TAILLE_MATRICE; i++ {
		for j := 0; j < TAILLE_MATRICE; j++ {
			if matrice[i][j] == 0 {
				fmt.Printf("       ")
			} else {
				fmt.Printf("%.2f ", matrice[i][j])
			}
		}

		fmt.Printf("\n")
	}
}

func TrouverVoisins(
	matrice [TAILLE_MATRICE][TAILLE_MATRICE]float64,
	depart_x, depart_y uint,
) [TAILLE_MATRICE][TAILLE_MATRICE]float64 {
	return trouverVoisinsRecursif(matrice, creerMatriceRes(), creerMatriceVisite(), depart_x, depart_y)
}

func trouverVoisinsRecursif(
	matrice [TAILLE_MATRICE][TAILLE_MATRICE]float64,
	res [TAILLE_MATRICE][TAILLE_MATRICE]float64,
	visites [TAILLE_MATRICE][TAILLE_MATRICE]bool,
	depart_x, depart_y uint,
) [TAILLE_MATRICE][TAILLE_MATRICE]float64 {
	x := depart_x
	y := depart_y

	if matrice[x][y] == 0 {
		return res
	}

	visites[x][y] = true
	res[x][y] = matrice[x][y]

	if x > 0 && !visites[x-1][y] && matrice[x-1][y] != 0 {
		res = trouverVoisinsRecursif(matrice, res, visites, x-1, y)
	}

	if x < uint(TAILLE_MATRICE-1) && !visites[x+1][y] && matrice[x+1][y] != 0 {
		res = trouverVoisinsRecursif(matrice, res, visites, x+1, y)
	}

	if y > 0 && !visites[x][y-1] && matrice[x][y-1] != 0 {
		res = trouverVoisinsRecursif(matrice, res, visites, x, y-1)
	}
	if y < uint(TAILLE_MATRICE-1) && !visites[x][y+1] && matrice[x][y+1] != 0 {
		res = trouverVoisinsRecursif(matrice, res, visites, x, y+1)
	}

	return res
}

func creerMatriceRes() (matrice [TAILLE_MATRICE][TAILLE_MATRICE]float64) {
	for i := 0; i < TAILLE_MATRICE; i++ {
		for j := 0; j < TAILLE_MATRICE; j++ {
			matrice[i][j] = 0
		}
	}

	return
}

func creerMatriceVisite() (matrice [TAILLE_MATRICE][TAILLE_MATRICE]bool) {
	for i := 0; i < TAILLE_MATRICE; i++ {
		for j := 0; j < TAILLE_MATRICE; j++ {
			matrice[i][j] = false
		}
	}

	return
}
