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

const N_GOROUTINE = 1

func Start() {
	// Handler for the main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { mainHandler(w, r) })
	http.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) { pointsHandler(w, r) })
	err := http.ListenAndServe(":8026", nil)
	fmt.Println("%v", err)
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
		Accuracy int
	}

	var jsonData JsonString
	err := json.Unmarshal(body, &jsonData)

	if err != nil {
		fmt.Println("Error decoding json")
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"bad json formatting\"}"))

		return
	}

	// Accuracy verification
	var accuracy gpsfiles.MapAccuracy

	switch jsonData.Accuracy {
	case 1:
		accuracy = gpsfiles.ACCURACY_1
	case 5:
		accuracy = gpsfiles.ACCURACY_5
	case 25:
		accuracy = gpsfiles.ACCURACY_25
	default:
		fmt.Println("Missing accuracy")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"missing accuracy\"}"))

		return
	}

	fmt.Println(accuracy)

	var squaresToShow []tileutils.Wgs84Square

	tiles := tileutils.ComputeTiles(jsonData.Lng, jsonData.Lat, jsonData.Deniv, accuracy, N_GOROUTINE)

	if len(tiles) == 0 {
		fmt.Println("Couldn't compute, probably invalid coordinates")
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("{\"error\": \"invalid coordinates\"}"))

		return
	}

	fmt.Printf("Computed %d tiles\n", len(tiles))

	for _, tile := range tiles {
		squaresToShow = append(squaresToShow, tile.ComputeOptimizedSquaresWgs()...)
	}

	var response = struct {
		TileSize int
		Tiles    []tileutils.Wgs84Square
	}{
		jsonData.Accuracy,
		squaresToShow,
	}

	w.Header().Set("Content-Type", "application/json") // On indique que la réponse sera au format JSON
	json.NewEncoder(w).Encode(response)                // On encode la liste des points en JSON et on l'envoie en réponse
}
