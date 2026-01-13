package algo

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// prend en argument le nom du fichier de BD de départ et les coordonnées x et y du point de départ en lambert et renvoie les indices i et j de la case de départ
func CaseDepart(xLambert, yLambert float64, cheminFichier string) (xMatrix, yMatrix int) {
	donnees, err := os.ReadFile(cheminFichier) // lire le fichier en question, data est en byte

	if err != nil {
		fmt.Println("Erreur lecture base de données:", err)
		return -1, -1
	}

	// fmt.Println(string(donnees)) // conversion de byte en string

	elem := strings.Fields(string(donnees)) // on sépare chaque element du fichier (sépare dès qu'il y a un espace ou retour à la ligne)

	// les coordonnées en lambert de la case d'indice [1000][0]
	xLL, err1 := strconv.ParseFloat(elem[5], 64) // conversion de string en float64
	yLL, err2 := strconv.ParseFloat(elem[7], 64) // conversion de string en float64

	if err1 != nil || err2 != nil {
		fmt.Println("Erreur lors de la récupération de des coordonnées", err)
		return -1, -1
	}


	xMatrix = int((xLambert - xLL)/25) // car entre chaque indice il y a 25m
	yMatrix = 1000-int((yLambert - yLL)/25)

	if xMatrix >= 1000 || yMatrix >= 1000 {
		fmt.Println("Coordonnées lambert en dehors de la case")
		return -1, -1
	}

	return
}
