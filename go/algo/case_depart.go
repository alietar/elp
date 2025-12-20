package algo

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// prend en argument le nom du fichier de BD de départ et les coordonnées x et y du point de départ en lambert et renvoie les indices i et j de la case de départ
func Case_depart(x0, y0 float64, fichier string) (i, j int) {

	donnees, err := os.ReadFile(fichier) // lire le fichier en question, data est en byte
	if err != nil {
		fmt.Println("Erreur lecture base de données:", err)
		return -1, -1
	}

	fmt.Println(string(donnees)) // conversion de byte en string

	elem := strings.Fields(string(donnees)) // on sépare chaque element du fichier (sépare dès qu'il y a un espace ou retour à la ligne)

	// les coordonnées en lambert de la case d'indice [1000][0]
	x, err1 := strconv.ParseFloat(elem[5], 64) // conversion de string en float64
	y, err2 := strconv.ParseFloat(elem[7], 64) // conversion de string en float64
	if err1 != nil || err2 != nil {
		fmt.Println("Erreur lors de la récupération de des coordonnées", err)
	}
	i = int(1000 - (x-x0)/25) // car entre chaque indice il y a 25m
	j = int(1000 + (y-y0)/25)

	return i, j

}
