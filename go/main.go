package main

import (
	"fmt"

	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/findfiles"
	// "github.com/alietar/elp/go/server"
)

func main() {
	longitude := "4.8723056"
	latitude := "45.7837778"

	folder := "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/"

	path, er := findfiles.GetFileForMyCoordinate(longitude, latitude, folder)
	path = folder + path

	if er != nil {
		fmt.Println(er)
		return
	}

	xLambert, yLambert, er := findfiles.FromGpsWgs84ToLambert93(longitude, latitude)

	if er != nil {
		fmt.Println(er)
		return
	}

	fmt.Printf("xLambert: %f, yLambert: %f\n", xLambert, yLambert)

	x, y := algo.CaseDepart(xLambert, yLambert, path)

	if x == -1 || y == -1 {
		fmt.Printf("Erreur sur le calcul de la case départ\n")
		return
	}

	fmt.Printf("x: %d, y: %d\n", x, y)
	mat := algo.CreationMatrice(path)
	mat2 := algo.PointsAtteignables(5, x, y, mat)

	fmt.Printf("Altitude : %f\n", mat[x][y])

	fullMatrix := algo.NewMatrix(1000)

	for i := range 1000 {
		for j := range 1000 {
			fullMatrix.Data[i][j] = mat2[i][j]
		}
	}

	fullMatrix = fullMatrix.FindNeighbors(x, y)

	smallMatrix, newX, newY := fullMatrix.Resize(x, y, 50)

	smallMatrix.ShowPrettyWithStart(newX, newY)
}
