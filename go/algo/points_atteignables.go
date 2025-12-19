package algo

// prend en argument la matrice des altitudes, les indices i0 et j0 de la case de départ et le dénivelé maximal d et renvoie la matrice contenant 0 pour les points non atteignables et leur altitudes pour les autres
func Points_atteignables(d, i0, j0 float64, Matrice [1000][1000]int) (Matrice_atteignable [1000][1000]int) {
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if Matrice[i][j] < abs(Matrice[i0][j0] - d){
				Matrice_atteignable[i][j] = Matrice[i][j]
			}
			else {
				Matrice_Matrice_atteignable[i][j] = 0
			}
		}
	}
	return Matrice_atteignable
}