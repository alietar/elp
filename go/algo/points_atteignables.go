package algo

import (
	"math"
)

// prend en argument la matrice des altitudes, les indices i0 et j0 de la case de départ et le dénivelé maximal d et renvoie la matrice contenant 0 pour les points non atteignables et leur altitudes pour les autres
func PointsAtteignables(d float64, i0, j0 int, matrice [1000][1000]float64) (matriceAtteignable [1000][1000]float64) {
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if math.Abs(matrice[i][j]-matrice[i0][j0]) < d {
				matriceAtteignable[i][j] = matrice[i][j]
			} else {
				matriceAtteignable[i][j] = 0
			}
		}
	}
	return matriceAtteignable
}
