package tileutils

import (
	"fmt"
	"math"
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

// startLongitude, startLatitude := *flag.String("long", "4.871928", "starting longitude"), *flag.String("lat", "45.7838052", "starting latitude")
func ComputeTiles(startLongitude, startLatitude, d float64) (returnTiles []*Tile) {
	// Getting the Lambert coordinates
	xLambert, yLambert, er := gpsfiles.ConvertWgs84ToLambert93(startLongitude, startLatitude)

	if er != nil {
		fmt.Println(er)
		return
	}

	fmt.Println("Starting coordinates")
	fmt.Printf(" - WGS84 => longitude: %f, latitude: %f\n", startLongitude, startLatitude)
	fmt.Printf(" - Lambert93 => xLambert: %f, yLambert: %f\n", xLambert, yLambert)

	// General initalization
	var exploredBorderPointsLambert [][2]float64

	adjacentTileCoordinatesChan := make(chan [2]float64, 100)
	doneTileMatricesChan := make(chan *Tile, 100)
	var wg sync.WaitGroup

	// Adding the first tile worker manually so the wait routine doesn't stop the program immediately
	wg.Add(1)
	adjacentTileCoordinatesChan <- [2]float64{xLambert, yLambert}
	startAlt := -1.0

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

		startAlt = addTileWorker(&wg, xLambert, yLambert, doneTileMatricesChan, adjacentTileCoordinatesChan, d, startAlt)
	}

	for tile := range doneTileMatricesChan {
		returnTiles = append(returnTiles, tile)
	}

	return

	// var pointsAffiche []algo.Coordonnee // Liste des points à afficher
	/*// Parcours de la matrice pour récupérer les coordonnées des points connectés
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
	}*/
}

func addTileWorker(wg *sync.WaitGroup, xLambert, yLambert float64, results chan *Tile, exploreAdj chan [2]float64, d, alt float64) float64 {
	fmt.Printf("New tile worker starting at: x=%f, y=%f\n", xLambert, yLambert)

	/////// !!!!! Use the same tile if the algo was already run on it
	tile, xStart, yStart := NewTileFromLambert(xLambert, yLambert)

	if alt == -1 {
		alt = tile.Altitudes[xStart][yStart]

		fmt.Printf("Starting at altitude %f\n", alt)
	}

	tile.CreatePotentiallyReachable(d, alt)

	if xStart == -1 || yStart == -1 {
		wg.Done()
		return alt
	}

	go FindNeighbors(tile, xStart, yStart, wg, results, exploreAdj)

	return alt
}

func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64, results chan *Tile) {
	wg.Wait()
	close(adjacentTileCoordinates)
	close(results)

	fmt.Println("All the tile workers finished")
}
