package algo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// On définit la structure attendue du JSON
type Depart struct {
	CodeDepartement string `json:"codeDepartement"`
}

// prend en paramètre les coordonnées x et y en wsg84 et retourne le département correspondant
func GetDepartement(x, y float64) int {
	url := fmt.Sprintf("https://geo.api.gouv.fr/communes?lat=%f&lon=%f&fields=codeDepartement", x, y)

	client := &http.Client{Timeout: 10 * time.Second} // Requête avec un Timeout
	rep, err := client.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return -1
	}
	defer rep.Body.Close() //rep.Body est le flux de données

	var data []Depart                             // car le JSON retourné est une liste
	err = json.NewDecoder(rep.Body).Decode(&data) // On décode le JSON qui est de la forme [{codeDepartement: "XX"}] en bytes
	if err != nil {
		fmt.Println("Erreur de décodage :", err)
		return -1
	}
	num, err := strconv.Atoi(data[0].CodeDepartement)
	if err != nil {
		fmt.Println("Erreur de conversion en entier :", err)
		return -1
	}
	return num
}
