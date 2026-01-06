package findfiles

import (
	"fmt"
	"testing"
)

func TestFromGpsWgs84ToLambert93(t *testing.T) {
	longitude := "2.3522"
	latitude := "48.8566"

	x, y, err := FromGpsWgs84ToLambert93(longitude, latitude)
	if err != nil {
		fmt.Println("Erreur de conversion :", err)
		return
	}

	fmt.Printf("Lambert-93 → X = %.2f m | Y = %.2f m\n", x, y)
}

func TestGetFilesNameFolder(t *testing.T) {
	folderPath := "/home/leopold/Documents/3TC/ELP/GO/elp/go/bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069"

	files, err := GetFilesNameFolder(folderPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier :", err)
		return
	}

	fmt.Println("Fichiers trouvés :")
	for _, file := range files {
		fmt.Println("-", file)
	}
}

func TestReadCoordinateLambert93File(t *testing.T) {
	filePath := "/home/leopold/Documents/3TC/ELP/GO/elp/go/bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/BDALTIV2_25M_FXX_0775_6550_MNT_LAMB93_IGN69.asc"

	xll, yll, cellsize, err := ReadCoordinateLambert93File(filePath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture des coordonnées :", err)
		return
	}

	fmt.Println("Coordonnées Lambert-93 trouvées :")
	fmt.Println("xllcorner :", xll)
	fmt.Println("yllcorner :", yll)
	fmt.Println("cell_size :", cellsize)
}

func TestGetFileForMyCoordinate(t *testing.T) {
	folderPath := "/home/leopold/Documents/3TC/ELP/GO/elp/go/bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069"

	longitude := "4.4483937"
	latitude  := "45.4922389"

	file, err := GetFileForMyCoordinate(longitude, latitude, folderPath)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Fichier contenant la coordonnée :")
	fmt.Println(file)
}


func TestBuildBDIndex(t *testing.T) {

	baseDir := "/home/leopold/Documents/3TC/ELP/GO/elp/go/bd"

	index, err := BuildBDIndex(baseDir)
	if err != nil {
		t.Fatalf("Erreur lors de la construction de l’index BD : %v", err)
	}

	if len(index) == 0 {
		t.Fatalf("Aucun fichier .asc trouvé dans %s", baseDir)
	}

	fmt.Printf("Nombre de fichiers indexés : %d\n", len(index))

	// Affichage des 5 premiers fichiers pour vérification visuelle
	limit := 5
	if len(index) < limit {
		limit = len(index)
	}

	for i := 0; i < limit; i++ {
		f := index[i]
		fmt.Printf(
			"%d) %s | xll=%.2f | yll=%.2f | cell=%.2f\n",
			i+1,
			f.Path,
			f.XllCorner,
			f.YllCorner,
			f.CellSize,
		)
	}
}


