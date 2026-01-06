package main

import (
	"encoding/json"
	"fmt"
<<<<<<< HEAD
	"sync"
	"flag"
	"slices"
=======
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
	"time"
>>>>>>> 5132283 (ajout de l'interface web interactive permettant d'afficher les points atteignables)

	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/findfiles"
	// "github.com/alietar/elp/go/server"
)

<<<<<<< HEAD

func main() {
	var exploredTilesPath []string
=======
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

func main() {
	longitude := "4.87226"
	latitude := "45.783844"
>>>>>>> 5132283 (ajout de l'interface web interactive permettant d'afficher les points atteignables)

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
	adjacentTileCoordinatesChan <- [2]float64{xLambert, yLambert}

	// This go routine stops the listenings for the channel adjacentTileCoordinates when no tile algorithm are at work
	go waitRoutine(&wg, adjacentTileCoordinatesChan)
	
	// addTileWorker(&wg, xLambert, yLambert, doneTileMatricesChan, adjacentTileCoordinatesChan, exploredTilesPath)

	

	// The for will wait for new adjacentTile until the waitRoutine closes the channel
	for tileCoordinates := range adjacentTileCoordinatesChan {
		folder := "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/"
		
		// Getting the right file
		path, er := findfiles.GetFileForMyCoordinate(tileCoordinates[0], tileCoordinates[1], folder)

		if er != nil {
			fmt.Println(er)
			continue
		}
		fmt.Printf("Tile is: %s\n", path)
		fmt.Printf("xLambert: %f, yLambert: %f\n", tileCoordinates[0], tileCoordinates[1])

		if slices.Contains(exploredTilesPath, path) {
			fmt.Printf("Tile at %s already explored\n", path)
			wg.Done()

			continue	
		}

		exploredTilesPath = append(exploredTilesPath, path)

		fmt.Println()
		fmt.Println(exploredTilesPath)
		fmt.Println()


		path = folder + path
		
		addTileWorker(&wg, tileCoordinates[0], tileCoordinates[1], doneTileMatricesChan, adjacentTileCoordinatesChan, path)
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


func addTileWorker(wg *sync.WaitGroup, xLambert, yLambert float64, results chan [][]float64, exploreAdj chan [2]float64, path string) {
	fmt.Printf("New tile worker starting at: x=%f, y=%f\n", xLambert, yLambert)

	mat := getMatrixFromLambert(xLambert, yLambert, path)

	go mat.FindNeighbors(mat.StartX, mat.StartY, wg, results, exploreAdj)
}


func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64) {
	wg.Wait()
	close(adjacentTileCoordinates)

	fmt.Println("All the tile workers finished")
}


func getMatrixFromLambert(xLambert, yLambert float64, path string) algo.Matrix {
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
		}
	}

<<<<<<< HEAD
	return fullMatrix
}
=======
	fmt.Printf("Altitude : %f\n", mat[x][y])

	fullMatrix = fullMatrix.FindNeighbors(x, y) // On trouve tous les points connectés à la case de départ

	xll, yll, _, _ := findfiles.ReadCoordinateLambert93File(path) // On lit les coordonnées en Lambert93 du coin inférieur gauche de la matrice

	var pointsAffiche []algo.Coordonnee // Liste des points à afficher

	// Parcours de la matrice pour récupérer les coordonnées des points connectés
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

	// Lance un serveur web
	http.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // On indique que la réponse sera au format JSON
		json.NewEncoder(w).Encode(pointsAffiche)           // On encode la liste des points en JSON et on l'envoie en réponse
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("carte.html") // On charge le fichier HTML
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// On prépare les données à injecter dans le template
		data := struct {
			Lat string
			Lng string
		}{
			Lat: latitude,
			Lng: longitude,
		}
		tmpl.Execute(w, data) // On exécute le template carte.html avec les données
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
>>>>>>> 5132283 (ajout de l'interface web interactive permettant d'afficher les points atteignables)
