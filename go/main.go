package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	// "encoding/json"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"slices"
	"sync"

	"github.com/alietar/elp/go/algo"

	"github.com/alietar/elp/go/findfiles"
	"github.com/alietar/elp/go/tileutils"
	// "github.com/alietar/elp/go/server"
)

func main() {
	// tileFolderPath := "./bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/"

	startLongitude, startLatitude := *flag.String("long", "4.871928", "starting longitude"), *flag.String("lat", "45.7838052", "starting latitude")

	// Getting the Lambert coordinates
	xLambert, yLambert, er := findfiles.FromGpsWgs84ToLambert93(startLongitude, startLatitude)

	if er != nil {
		fmt.Println(er)
		return
	}

	fmt.Println("Starting coordinates")
	fmt.Printf(" - WGS84 => longitude: %s, latitude: %s\n", startLongitude, startLatitude)
	fmt.Printf(" - Lambert93 => xLambert: %f, yLambert: %f\n", xLambert, yLambert)

	// General initalization
	var exploredBorderPointsLambert [][2]float64

	adjacentTileCoordinatesChan := make(chan [2]float64, 100)
	doneTileMatricesChan := make(chan *tileutils.Tile, 100)
	var wg sync.WaitGroup

	// Adding the first tile worker manually so the wait routine doesn't stop the program immediately
	wg.Add(1)
	adjacentTileCoordinatesChan <- [2]float64{xLambert, yLambert}

	// This go routine stops the listenings for the channel
	// adjacentTileCoordinates when no tile algorithm are at work
	go waitRoutine(&wg, adjacentTileCoordinatesChan, doneTileMatricesChan)

	// The "for" will wait for new adjacentTile until the waitRoutine closes the channel
	for entryPointCoordinates := range adjacentTileCoordinatesChan {
		xLambert := entryPointCoordinates[0]
		yLambert := entryPointCoordinates[1]

		////// !!!! CHECK Explorated coordinates only for the same tile path
		// Dont'go explore a tile if it was already explored starting from the same border point
		skip := false
		for _, point := range exploredBorderPointsLambert {
			if math.Sqrt(math.Pow((point[0]-entryPointCoordinates[0]), 2)+math.Pow((point[1]-entryPointCoordinates[1]), 2)) < 1000 {
				fmt.Printf("Tile already explored by this way\n")

				skip = true
				break
			}
		}

		if skip {
			wg.Done()
			continue
		}

		/*if slices.Contains(exploredBorderPointsLambert, entryPointCoordinates) {
			fmt.Printf("Tile already explored by this way\n")
			wg.Done()

			continue
		}*/
		fmt.Printf("xLambert: %f, yLambert: %f\n", xLambert, yLambert)

		exploredBorderPointsLambert = append(exploredBorderPointsLambert, entryPointCoordinates)

		addTileWorker(&wg, xLambert, yLambert, doneTileMatricesChan, adjacentTileCoordinatesChan)
	}

	var pointsAffiche []algo.Coordonnee // Liste des points à afficher

	for tile := range doneTileMatricesChan {
		// Parcours de la matrice pour récupérer les coordonnées des points connectés
		for i := 0; i < 1000; i++ {
			for j := 0; j < 1000; j++ {
				if tile.Reachable[i][j] == true {
					X := tile.XLambertLL + float64(i)*25.0             // Conversion de l'indice de colonne en coordonnée Lambert93
					Y := (tile.YLambertLL) + 25000.0 - float64(j)*25.0 // Conversion de l'indice de ligne en coordonnée Lambert93
					lat, lng, _ := algo.FromLambert93ToGpsWgs84(X, Y)  // Conversion en coordonnées GPS WGS84
					pt := algo.Coordonnee{Lat: lat, Lng: lng}

					if !slices.Contains(pointsAffiche, pt) {
						pointsAffiche = append(pointsAffiche, pt) // Ajout à la liste des points à afficher
					}
				}
			}
		}
	}

	httpServer(pointsAffiche)

	// This code is executes when the channel adjacentTileCoordinatesChan is closed
	/*
		var pointsAffiche []algo.Coordonnee // Liste des points à afficher

		for tile := range doneTileMatricesChan {
			xll, yll, _, _ := findfiles.ReadCoordinateLambert93File(path) // On lit les coordonnées en Lambert93 du coin inférieur gauche de la matrice

			// Parcours de la matrice pour récupérer les coordonnées des points connectés
			for i := 0; i < 1000; i++ {
				for j := 0; j < 1000; j++ {
					if tile.Reachable[i][j] == true {
						X := tile.XLambertLL + float64(j)*25.0                                                 // Conversion de l'indice de colonne en coordonnée Lambert93
						Y := (tile.YLambertLL) + 25000.0) - float64(i)*25.0                                     // Conversion de l'indice de ligne en coordonnée Lambert93
						lat, lng, _ := algo.FromLambert93ToGpsWgs84(X, Y)                          // Conversion en coordonnées GPS WGS84
						pointsAffiche = append(pointsAffiche, algo.Coordonnee{Lat: lat, Lng: lng}) // Ajout à la liste des points à afficher
					}
				}
			}
		}

		httpServer(pointsAffiche)*/
}

func addTileWorker(wg *sync.WaitGroup, xLambert, yLambert float64, results chan *tileutils.Tile, exploreAdj chan [2]float64) {
	fmt.Printf("New tile worker starting at: x=%f, y=%f\n", xLambert, yLambert)

	/////// !!!!! Use the same tile if the algo was already run on it
	tile, xStart, yStart := tileutils.NewTileFromLambert(xLambert, yLambert)

	if xStart == -1 || yStart == -1 {
		wg.Done()
		return
	}

	go tileutils.FindNeighbors(tile, xStart, yStart, wg, results, exploreAdj)
}

func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64, results chan *tileutils.Tile) {
	wg.Wait()
	close(adjacentTileCoordinates)
	close(results)

	fmt.Println("All the tile workers finished")
}

/*
	fmt.Printf("Altitude : %f\n", mat[x][y])

	fullMatrix = fullMatrix.FindNeighbors(x, y) // On trouve tous les points connectés à la case de départ

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

*/

func httpServer(pointsAffiche []algo.Coordonnee) {
	http.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pointsAffiche)
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
			Lat: "45.430121",
			Lng: "4.582209",
		}
		tmpl.Execute(w, data) // On exécute le template carte.html avec les données
	})

	fmt.Println("Calcul terminé. Préparation de la carte...")

	http.ListenAndServe(":8080", nil) // On démarre le serveur web sur le port 8080
}
