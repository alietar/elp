package algo

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"find_files"
)

func Creation_matrice() [1000][1000]float64 {

	var Matrice [1000][1000]float64

	fichier := find_files. // récupérer le resultat de la fonction qui trouve le bon fichier de la BD
	donnees, err := os.ReadFile(fichier) // lire le fichier en question, data est en byte
	if err != nil {
		fmt.Println("Erreur lecture base de données:", err)
		return Matrice
	}

	fmt.Println(string(donnees)) // conversion de byte en string

	elem := strings.Fields(string(donnees)) // on sépare chaque element du fichier (sépare dès qu'il y a un espace ou retour à la ligne)

	altitudes := elem[12:] // les altitudes démarrent au 13e élément

	k := 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			alt, err := strconv.ParseFloat(altitudes[k], 64) // conversion de string en float64
			if err != nil {
				fmt.Println("Erreur lors de la récupération de l'altitude", err)
			}
			Matrice[i][j] = alt
			k += 1
		}
	}

	return Matrice
}
