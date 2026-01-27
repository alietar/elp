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

	// Concurrent tools init
	tileCache := NewTileCache()
	adjacentTileCoordinatesChan := make(chan [2]float64, 20000)
	// adjacentTileCoordinatesChan := make(chan [2]float64, int(10000.0/float64(nWorker)))
	var wg sync.WaitGroup

	// Loading the first tile manually to find starting altitude
	startTile, xStart, yStart := tileCache.GetOrLoad(xLambert, yLambert, accuracy)

	// Did not find the tile, quitting
	if xStart == -1 || yStart == -1 {
		wg.Done()
		return nil
	}

	startAlt := startTile.Altitudes[xStart][yStart]

	wg.Add(1)
	adjacentTileCoordinatesChan <- [2]float64{xLambert, yLambert}

	for i := range nWorker {
		fmt.Printf("Launching worker %d\n", i)
		go tileWorker(&wg, tileCache, adjacentTileCoordinatesChan, d, startAlt, accuracy)
	}

	// This go routine stops the listenings for the channel
	// adjacentTileCoordinates when no worker are at work
	// Which then stops the workers
	// go waitRoutine(&wg, adjacentTileCoordinatesChan)

	wg.Wait()
	close(adjacentTileCoordinatesChan)

	fmt.Println("All the tile workers finished")

	return tileCache.GetValuesSlice()
}

func tileWorker(wg *sync.WaitGroup,
	tileCache *TileCache,
	exploreChannel chan [2]float64,
	d, alt float64,
	accuracy gpsfiles.MapAccuracy,
) {
	// fmt.Println("Hey I'm a worker")
	for entryPointCoordinates := range exploreChannel {
		// fmt.Println("Received coordinates")
		xLambert := entryPointCoordinates[0]
		yLambert := entryPointCoordinates[1]

		workerComputeCoordinates(wg, tileCache, exploreChannel, xLambert, yLambert, d, alt, accuracy)
	}
}

func workerComputeCoordinates(wg *sync.WaitGroup,
	tileCache *TileCache,
	exploreChannel chan [2]float64,
	xLambert, yLambert, d, alt float64,
	accuracy gpsfiles.MapAccuracy,
) (skipped bool) {
	tile, xStart, yStart := tileCache.GetOrLoad(xLambert, yLambert, accuracy)

	// Did not find the tile, skipping
	if xStart == -1 || yStart == -1 {
		wg.Done()
		return true
	}

	tile.Mutex.Lock()

	// Init matrices if new tile
	if tile.PotentiallyReachable == nil {
		tile.CreatePotentiallyReachable(d, alt)
	}

	// Skip if already reachable
	if tile.Reachable[xStart][yStart] {
		tile.Mutex.Unlock()
		wg.Done()
		return true
	}

	// Writing tile here to prevent other goroutines from starting
	tile.Reachable[xStart][yStart] = true
	tile.Mutex.Unlock()

	FindNeighbors(tile, xStart, yStart, wg, exploreChannel)

	return false
}

func waitRoutine(wg *sync.WaitGroup, adjacentTileCoordinates chan [2]float64) {
}
