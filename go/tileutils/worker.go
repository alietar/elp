package tileutils

import (
	"fmt"
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

func ComputeTiles(startLongitude, startLatitude, d float64, accuracy gpsfiles.MapAccuracy) (returnTiles []*Tile) {
	// Getting the Lambert coordinates
	xLambert, yLambert, er := gpsfiles.ConvertWgs84ToLambert93(startLongitude, startLatitude)

	if er != nil {
		fmt.Println(er)
		return
	}

	fmt.Println("Starting coordinates")
	fmt.Printf(" - WGS84 => longitude: %f, latitude: %f\n", startLongitude, startLatitude)
	fmt.Printf(" - Lambert93 => xLambert: %f, yLambert: %f\n", xLambert, yLambert)

	tileCache := NewTileCache()

	adjacentTileCoordinatesChan := make(chan [2]float64, 5000)
	var wg sync.WaitGroup

	// Adding the first tile worker manually so the wait routine doesn't stop the program immediately
	wg.Add(1)
	adjacentTileCoordinatesChan <- [2]float64{xLambert, yLambert}
	startAlt := -1.0

	// This go routine stops the listenings for the channel
	// adjacentTileCoordinates when no tile algorithm are at work
	go waitRoutine(&wg, adjacentTileCoordinatesChan)

	for entryPointCoordinates := range adjacentTileCoordinatesChan {
		xLam := entryPointCoordinates[0]
		yLam := entryPointCoordinates[1]

		// 1. Récupérer la tuile depuis le cache
		tile, xStart, yStart := tileCache.GetOrLoad(xLam, yLam, accuracy)

		if xStart == -1 || yStart == -1 {
			wg.Done()
			continue
		}

		// 2. VERROUILLAGE CRITIQUE
		tile.Mutex.Lock()

		// Init lazy de la matrice booléenne
		if tile.PotentiallyReachable == nil {
			if startAlt == -1 {
				startAlt = tile.Altitudes[xStart][yStart]
			}
			tile.CreatePotentiallyReachable(d, startAlt)
		}

		// 3. LE TEST DE L'AMNÉSIE
		// Si la case est déjà cochée, ON ARRÊTE TOUT DE SUITE.
		// On ne lance pas de worker, on ne fait pas de wg.Add.
		if tile.Reachable[xStart][yStart] {
			tile.Mutex.Unlock()
			wg.Done()
			continue
		}

		// 4. LE MARQUAGE PRÉVENTIF (Atomic Set)
		// On coche la case MAINTENANT pour que si un autre worker arrive
		// une milliseconde plus tard, il voit que c'est pris.
		tile.Reachable[xStart][yStart] = true

		tile.Mutex.Unlock()

		// 5. Lancement du worker
		// Note: On passe 'tile' directement, plus besoin de recharger le fichier dedans
		go FindNeighbors(tile, xStart, yStart, &wg, adjacentTileCoordinatesChan)
	}

	return tileCache.GetValuesSlice()
}

// Version ajustée de addTileWorker pour aller avec le code ci-dessus
func addTileWorker(wg *sync.WaitGroup,
	tile *Tile, // On passe la tuile directement
	xStart, yStart int, // Et les indices calculés
	d, alt float64,
	exploreAdj chan [2]float64,
) float64 {

	// Init thread-safe de la tuile si nécessaire
	tile.Mutex.Lock()
	if tile.PotentiallyReachable == nil {
		if alt == -1 {
			alt = tile.Altitudes[xStart][yStart]
		}
		tile.CreatePotentiallyReachable(d, alt)
	}
	// Petite double sécurité (si un autre worker est passé entre le check du main et ici)
	if tile.Reachable[xStart][yStart] {
		tile.Mutex.Unlock()
		return alt
	}
	tile.Mutex.Unlock()

	// C'est SEULEMENT ICI qu'on incrémente le WaitGroup
	wg.Add(1)
	// On lance la version Itérative (non récursive)
	go FindNeighbors(tile, xStart, yStart, wg, exploreAdj)

	return alt
}

func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64) {
	wg.Wait()
	close(adjacentTileCoordinates)

	fmt.Println("All the tile workers finished")
}
