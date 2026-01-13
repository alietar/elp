package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"flag"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/tile"

	"github.com/alietar/elp/go/findfiles"
	// "github.com/alietar/elp/go/server"
)


func main() {
	longPtr := flag.String("long", "4.582209", "starting longitude")
	latPtr := flag.String("lat", "45.430121", "starting latitude")

	startLongitude := *longPtr
	startLatitude := *latPtr
	
	fmt.Printf("longitude: %s, latitude: %s\n", startLongitude, startLatitude)


	// Getting the Lambert coordinates
	xLambert, yLambert, er := findfiles.FromGpsWgs84ToLambert93(startLongitude, startLatitude)

	if er != nil {
		fmt.Println(er)
		return
	}

	fmt.Printf("xLambert: %f, yLambert: %f\n", xLambert, yLambert)


	// Channels and waitgroup initalization
	adjacentTileCoordinatesChan := make(chan [2]float64, 100)
	doneTileMatricesChan:= make(chan [][]float64, 100)
	var wg sync.WaitGroup


	// Adding the first tile worker manually so the wait routine doesn't stop the program immediately
	wg.Add(1)

	// This go routine stops the listenings for the channel adjacentTileCoordinates when no tile algorithm are at work
	go waitRoutine(&wg, adjacentTileCoordinatesChan)
	
	addTileWorker(&wg, xLambert, yLambert, doneTileMatricesChan, adjacentTileCoordinatesChan)
	

	// The for will wait for new adjacentTile until the waitRoutine closes the channel
	for tileCoordinates := range adjacentTileCoordinatesChan {
		addTileWorker(&wg, tileCoordinates[0], tileCoordinates[1], doneTileMatricesChan, adjacentTileCoordinatesChan)
	}


	// This code is executes when the channel adjacentTileCoordinatesChan is closed
	for tileMatriceResult := range doneTileMatricesChan {
		/*for i := range 1000 {
			for j := range 1000 {
				fmt.Printf("%f ", tileMatriceResult[i][j])
			}

			fmt.Print("\n")
		}*/

		fmt.Println(len(tileMatriceResult))
	}
	

	
	/*fmt.Println(<-adjacentTileCoordinates)
doneTilePtr<-[][]float64{{1, 2}}
	fmt.Println(<-doneTilePtr)

	wg.Done()

	time.Sleep(1 * time.Second)
	adjacentTileCoordinates <- [2]float64{xLambert, yLambert}*/
	
	
	/*

	// Getting the points in a matrix struct
	startMatrix := getMatrixFromLambert(xLambert, yLambert)

	var results := make([][][]float64, 100)

	threadNumber := 2

    resultChannel := make(chan [][]int, threadNumber)


	wg.Add(1)
	startMatrix = go startMatrix.FindNeighbors(startMatrix.startX, startMatrix.startY, &wg, resultChannel)

	go func() {
		wg.Wait()
		close(results)
	}

	count := 0
	for res := range resultChannel {
		results[count] := make([][]float64, 1000)
		for i := range 1000 {
			results[count][i] := make([]int, 1000)
			
			for j := range 1000 {
				results[count][i][j] = res[i][j]
			}
		}

		count++
	}

	fmt.Printf("Terminé. %d cellules remplies.\n", count)




	smallMatrix, newX, newY := fullMatrix.Resize(x, y, 50)

	smallMatrix.ShowPrettyWithStart(newX, newY)*/
}


func addTileWorker(wg *sync.WaitGroup, xLambert, yLambert float64, results chan [][]float64, exploreAdj chan [2]float64) {

	fmt.Printf("New tile worker starting at: x=%f, y=%f\n", xLambert, yLambert)

	mat := getMatrixFromLambert(xLambert, yLambert)

	go mat.FindNeighbors(mat.StartX, mat.StartY, wg, results, exploreAdj)
}


func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64) {
	wg.Wait()
	close(adjacentTileCoordinates)

	fmt.Println("All the tile workers finished")
}


func getMatrixFromLambert(xLambert, yLambert float64) algo.Matrix {
	folder := "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/"
	
	// Getting the right file
	path, er := findfiles.GetFileForMyCoordinate(xLambert, yLambert, folder)

	fmt.Printf("Path to file: %s\n", path)

	path = folder + path

	if er != nil {
		fmt.Println(er)
		return algo.Matrix{}
	}

	x, y := algo.CaseDepart(xLambert, yLambert, path)

	if x == -1 || y == -1 {
		fmt.Printf("Erreur sur le calcul de la case départ\n")
		return algo.Matrix{}
	}

	fmt.Printf("x: %d, y: %d\n", x, y)
	mat := algo.CreationMatrice(path)
	mat = algo.PointsAtteignables(10, x, y, mat)

	fmt.Printf("Altitude : %f\n", mat[x][y])

	fullMatrix := algo.NewMatrix(1000)
	fullMatrix.LambertX = xLambert
	fullMatrix.LambertY = yLambert
	fullMatrix.StartX = x
	fullMatrix.StartY = y

	for i := range 1000 {
		for j := range 1000 {
			fullMatrix.Data[i][j] = mat[i][j]
	// Handler pour la page principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		latitude := strings.TrimSpace(r.URL.Query().Get("lat")) // Récupère les paramètres lat, lng et deniv de l'URL, TrimSpace nettoie les espaces inutiles
		longitude := strings.TrimSpace(r.URL.Query().Get("lng"))
		denivele_srtg := strings.TrimSpace(r.URL.Query().Get("deniv"))

		// On prépare les données à injecter dans le template
		data := struct {
			Lat string
			Lng string
			Size int
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

						
			fullMatrixOptimize := tile.OptimizeSquares(fullMatrix.Data,xll,yll,25)


			// Parcours de la matrice pour récupérer les coordonnées des points connectés à afficher
			for i := 0; i < 1000; i++ {
				for j := 0; j < 1000; j++ {
					if fullMatrixOptimize.Data[i][j] != 0 {
						//X := xll + float64(j)*25.0                                                 // Conversion de l'indice de colonne en coordonnée Lambert93
						//Y := (yll + 25000.0) - float64(i)*25.0                                     // Conversion de l'indice de ligne en coordonnée Lambert93
						
						X = fullMatrixOptimize.CenterX
						Y = fullMatrixOptimize.CenterY
						size = fullMatrixOptimize.Size

						lat, lng, _ := algo.FromLambert93ToGpsWgs84(X, Y)                          // Conversion en coordonnées GPS WGS84
						pointsAffiche = append(pointsAffiche, algo.Coordonnee{Lat: lat, Lng: lng, Size: size}) // Ajout à la liste des points à afficher
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
