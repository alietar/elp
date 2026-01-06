package main

import (
	"fmt"
	"sync"
	"flag"
	"slices"

	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/findfiles"
	// "github.com/alietar/elp/go/server"
)


func main() {
	var exploredTilesPath []string

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

	return fullMatrix
}