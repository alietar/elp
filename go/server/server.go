package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/alietar/elp/go/tileutils"
)

func Start() {
	// Handler for the main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { mainHandler(w, r) })
	http.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) { pointsHandler(w, r) })
	http.ListenAndServe(":8080", nil) // Starts the server
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	tmpl, err := template.ParseFiles("server/map.html") // On charge le fichier HTML
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil) // On exécute le template carte.html avec les données
}

func pointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	////// Add error management like point out of map etc...
	fmt.Println("New request for points")
	body, _ := io.ReadAll(r.Body)
	fmt.Println(string(body))

	type JsonString struct {
		Lat   float64
		Lng   float64
		Deniv float64
	}
	var jsonData JsonString
	json.Unmarshal(body, &jsonData)

	// Gets the query parameters in the URL
	/*
		latitudeStr := strings.TrimSpace(r.URL.Query().Get("lat")) // TrimSpace remove useless spaces
		longitudeStr := strings.TrimSpace(r.URL.Query().Get("lng"))
		elevationStr := strings.TrimSpace(r.URL.Query().Get("deniv"))

		latitude, err1 := strconv.ParseFloat(latitudeStr, 64)
		longitude, err2 := strconv.ParseFloat(longitudeStr, 64)
		elevation, err3 := strconv.ParseFloat(elevationStr, 64)

		// Verify that query parameters are correct
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("Error in the parsing of the query parameters")

			fmt.Println(err1)
		}
	*/

	if jsonData.Lat > 0 && jsonData.Lng > 0 && jsonData.Deniv > 0 {
		// pointsAffiche = nil // Réinitialise la liste des points à afficher à chaque requête
		var squaresToShow []tileutils.Wgs84Square

		tiles := tileutils.ComputeTiles(jsonData.Lng, jsonData.Lat, jsonData.Deniv)

		for _, tile := range tiles {
			squaresToShow = append(squaresToShow, tile.ComputeOptimizedSquaresWgs()...)
		}

		w.Header().Set("Content-Type", "application/json") // On indique que la réponse sera au format JSON
		json.NewEncoder(w).Encode(squaresToShow)           // On encode la liste des points en JSON et on l'envoie en réponse
	}

	// json.NewEncoder(w).Encode(squaresToShow)           // On encode la liste des points en JSON et on l'envoie en réponse
}
