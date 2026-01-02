package main

import (
	"fmt"

	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/findfiles"
	// "github.com/alietar/elp/go/server"
)

func main() {
	longitude := "3.6554553"
	latitude := "46.0416860"

	xLambert, yLambert, er := findfiles.FromGpsWgs84ToLambert93(longitude, latitude)
	fmt.Printf("xLambert: %f, yLambert: %f\n", xLambert, yLambert)
	fmt.Println(er)

	// xLambert = 774987.500000000000 - 24300
	// yLambert = 6525012.500000000000 + 24300

	x, y := algo.CaseDepart(xLambert, yLambert, "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/BDALTIV2_25M_FXX_0775_6550_MNT_LAMB93_IGN69.asc")

	fmt.Printf("x: %d, y: %d\n", x, y)

	mat := algo.CreationMatrice("./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/BDALTIV2_25M_FXX_0775_6550_MNT_LAMB93_IGN69.asc")
	mat2 := algo.PointsAtteignables(2, x, y, mat)

	fmt.Printf("%f\n", mat[x][y])

	testSize := 60

	matrix := algo.NewMatrix(testSize)

	for i := range testSize {
		for j := range testSize {
			matrix.Data[i][j] = mat2[i][j]
		}
	}

	matrix.ShowPrettyWithStart(x, y)

	matrix2 := matrix.FindNeighbors(x, y)
	matrix2.ShowPrettyWithStart(x, y)
}
