package server

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"runtime"

	"github.com/alietar/elp/go/gpsfiles"
	"github.com/alietar/elp/go/tileutils"
)

func Start() {
	http.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) { pointsHandler(w, r) })

	fs := http.FileServer(http.Dir("../elm/"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":8026", nil)
	fmt.Println("%v", err)
}

func pointsHandler(w http.ResponseWriter, r *http.Request) {
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

	nCPU := runtime.NumCPU()
	nExploreWorker := int(math.Sqrt(float64(nCPU)))
	nFileWorker := int(nCPU / nExploreWorker)

	tiles := tileutils.ComputeTiles(jsonData.Lng, jsonData.Lat, jsonData.Deniv, accuracy, nExploreWorker, nFileWorker)

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
