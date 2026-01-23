package tileutils

import (
	"fmt"
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

// startLongitude, startLatitude := *flag.String("long", "4.871928", "starting longitude"), *flag.String("lat", "45.7838052", "starting latitude")
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

		// path, _, _, _ := gpsfiles.GetFileForMyCoordinate(xLambert, yLambert, "./db/"+string(accuracy)+"/")

		// fmt.Printf("xLambert: %f, yLambert: %f\n", xLambert, yLambert)

		startAlt = addTileWorker(
			xLambert, yLambert, d, startAlt,
			accuracy,
			tileCache,
			&wg,
			doneTileMatricesChan,
			adjacentTileCoordinatesChan,
		)
	}

	uniqueTiles := make(map[*Tile]bool)
	for tile := range doneTileMatricesChan {
		if !uniqueTiles[tile] {
			returnTiles = append(returnTiles, tile)
			uniqueTiles[tile] = true
		}
	}

	return

	// for tile := range doneTileMatricesChan {
	// returnTiles = append(returnTiles, tile)
	// }

	// return
}

func addTileWorker(
	xLambert, yLambert, d, alt float64,
	accuracy gpsfiles.MapAccuracy,
	tc *TileCache,
	wg *sync.WaitGroup,
	results chan *Tile,
	exploreAdj chan [2]float64,
) float64 {

	tile, xStart, yStart, _ := tc.GetOrLoad(xLambert, yLambert, accuracy)

	/////// !!!!! Use the same tile if the algo was already run on it
	// tile, xStart, yStart := NewTileFromLambert(xLambert, yLambert, accuracy)

	if xStart == -1 || yStart == -1 {
		wg.Done()
		return alt
	}

	//

	tile.Mutex.Lock()
	if tile.PotentiallyReachable == nil { // Null if it is a new tile
		// Si alt == -1 (cas du tout premier point), on le set
		if alt == -1 {
			alt = tile.Altitudes[xStart][yStart]
			fmt.Printf("Starting global altitude set to %f\n", alt)
		}
		tile.CreatePotentiallyReachable(d, alt)
	}

	// 3. VÉRIFICATION CRITIQUE : Est-ce qu'on est déjà passé par ce pixel précis ?
	if tile.Reachable[xStart][yStart] {
		// On a déjà exploré ce point précis via un autre chemin.
		// On arrête là pour cette branche.
		tile.Mutex.Unlock()
		wg.Done()
		return alt
	}

	fmt.Printf("New tile worker starting at: x=%f, y=%f\n", xLambert, yLambert)

	tile.Mutex.Unlock()

	go FindNeighbors(tile, xStart, yStart, wg, results, exploreAdj)

	return alt
}

func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64, results chan *Tile) {
	wg.Wait()
	close(adjacentTileCoordinates)
	close(results)

	fmt.Println("All the tile workers finished")
}
