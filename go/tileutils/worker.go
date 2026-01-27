package tileutils

import (
	"fmt"
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

func ComputeTiles(startLongitude, startLatitude, d float64, accuracy gpsfiles.MapAccuracy, nWorker int) (returnTiles []*Tile) {
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

	adjacentTileCoordinatesChan := make(chan [2]float64, nWorker)
	var wg sync.WaitGroup

	// Adding the first tile worker manually so the wait routine doesn't stop the program immediately
	wg.Add(1)
	adjacentTileCoordinatesChan <- [2]float64{xLambert, yLambert}
	startAlt := -1.0

	// This go routine stops the listenings for the channel
	// adjacentTileCoordinates when no tile algorithm are at work
	go waitRoutine(&wg, adjacentTileCoordinatesChan)

	for entryPointCoordinates := range adjacentTileCoordinatesChan {
		xLambert := entryPointCoordinates[0]
		yLambert := entryPointCoordinates[1]

		startAlt = addTileWorker(&wg, tileCache, adjacentTileCoordinatesChan, xLambert, yLambert, d, startAlt, accuracy)
	}

	return tileCache.GetValuesSlice()
}

func addTileWorker(wg *sync.WaitGroup,
	tileCache *TileCache,
	exploreAdj chan [2]float64,
	xLambert, yLambert, d, alt float64,
	accuracy gpsfiles.MapAccuracy,
) float64 {
	tile, xStart, yStart := tileCache.GetOrLoad(xLambert, yLambert, accuracy)

	// Did not find the tile, skipping
	if xStart == -1 || yStart == -1 {
		wg.Done()
		return alt
	}

	tile.Mutex.Lock()

	if tile.PotentiallyReachable == nil {
		if alt == -1 {
			alt = tile.Altitudes[xStart][yStart]
		}

		tile.CreatePotentiallyReachable(d, alt)
	}

	if tile.Reachable[xStart][yStart] {
		tile.Mutex.Unlock()
		wg.Done()

		return alt
	}

	// Writing tile hear to prevent other goroutines from starting
	tile.Reachable[xStart][yStart] = true

	tile.Mutex.Unlock()

	go FindNeighbors(tile, xStart, yStart, wg, exploreAdj)

	return alt
}

func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64) {
	wg.Wait()
	close(adjacentTileCoordinates)

	fmt.Println("All the tile workers finished")
}
