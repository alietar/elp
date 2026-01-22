package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/alietar/elp/go/gpsfiles"
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
		Lat      float64
		Lng      float64
		Deniv    float64
		Accuracy string
	}

	var jsonData JsonString
	json.Unmarshal(body, &jsonData)

	// Accuracy verification
	var accuracy gpsfiles.MapAccuracy

	switch jsonData.Accuracy {
	case "1":
		accuracy = gpsfiles.ACCURACY_1
	case "5":
		accuracy = gpsfiles.ACCURACY_1
	case "25":
		accuracy = gpsfiles.ACCURACY_1
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"missing accuracy\"}"))

		return
	}

	fmt.Println(accuracy)

	var squaresToShow []tileutils.Wgs84Square

	tiles := tileutils.ComputeTiles(jsonData.Lng, jsonData.Lat, jsonData.Deniv)

	if len(tiles) == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("{\"error\": \"invalid coordinates\"}"))

		return
	}

	fmt.Printf("Computed %d tiles\n", len(tiles))

	for _, tile := range tiles {
		squaresToShow = append(squaresToShow, tile.ComputeOptimizedSquaresWgs()...)
	}

	w.Header().Set("Content-Type", "application/json") // On indique que la réponse sera au format JSON
	json.NewEncoder(w).Encode(squaresToShow)           // On encode la liste des points en JSON et on l'envoie en réponse
}
