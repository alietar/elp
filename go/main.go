package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/findfiles"
	// "github.com/alietar/elp/go/server"
)

// openBrowser sert à détecter le système d'exploitation et ouvrir le navigateur par défaut avec l'URL donnée
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du navigateur:", err)
	}
}

var pointsAffiche []algo.Coordonnee // Liste des points à afficher, global pour être accessible dans les handlers

func main() {

	folder := "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/"

	// Handler pour la page principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		latitude := strings.TrimSpace(r.URL.Query().Get("lat")) // Récupère les paramètres lat, lng et deniv de l'URL, TrimSpace nettoie les espaces inutiles
		longitude := strings.TrimSpace(r.URL.Query().Get("lng"))
		denivele_srtg := strings.TrimSpace(r.URL.Query().Get("deniv"))

		// On prépare les données à injecter dans le template
		data := struct {
			Lat string
			Lng string
		}{
			Lat: latitude,
			Lng: longitude,
		}

		if latitude != "" && longitude != "" && denivele_srtg != "" {
			pointsAffiche = nil // Réinitialise la liste des points à afficher à chaque requête

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

			x, y := algo.CaseDepart(xLambert, yLambert, path)

			if x == -1 || y == -1 {
				fmt.Printf("Erreur sur le calcul de la case départ\n")
				return
			}

			denivele, err := strconv.ParseFloat(denivele_srtg, 64)
			if err != nil {
				fmt.Println(err)
				return
			}

			mat := algo.CreationMatrice(path)
			mat2 := algo.PointsAtteignables(denivele, x, y, mat)

			fullMatrix := algo.NewMatrix(1000)

			for i := range 1000 {
				for j := range 1000 {
					fullMatrix.Data[i][j] = mat2[i][j]
				}
			}

			fullMatrix = fullMatrix.FindNeighbors(x, y) // On trouve tous les points connectés à la case de départ

			xll, yll, _, _ := findfiles.ReadCoordinateLambert93File(path) // On lit les coordonnées en Lambert93 du coin inférieur gauche de la matrice

			// Parcours de la matrice pour récupérer les coordonnées des points connectés à afficher
			for i := 0; i < 1000; i++ {
				for j := 0; j < 1000; j++ {
					if fullMatrix.Data[i][j] != 0 {
						X := xll + float64(j)*25.0                                                 // Conversion de l'indice de colonne en coordonnée Lambert93
						Y := (yll + 25000.0) - float64(i)*25.0                                     // Conversion de l'indice de ligne en coordonnée Lambert93
						lat, lng, _ := algo.FromLambert93ToGpsWgs84(X, Y)                          // Conversion en coordonnées GPS WGS84
						pointsAffiche = append(pointsAffiche, algo.Coordonnee{Lat: lat, Lng: lng}) // Ajout à la liste des points à afficher
					}
				}
			}
		}

		tmpl, err := template.ParseFiles("carte.html") // On charge le fichier HTML
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data) // On exécute le template carte.html avec les données
	})

	// Lance un serveur web
	http.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // On indique que la réponse sera au format JSON
		json.NewEncoder(w).Encode(pointsAffiche)           // On encode la liste des points en JSON et on l'envoie en réponse
	})

	fmt.Println("Calcul terminé. Préparation de la carte...")

	// On lance l'ouverture automatique dans un fil séparé
	go func() {
		// On attend 1 seconde pour être sûr que le serveur est bien prêt
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:8080") // On ouvre le navigateur par défaut avec l'URL du serveur
	}()

	http.ListenAndServe(":8080", nil) // On démarre le serveur web sur le port 8080

}
