package tileutils

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

func NewTileFromLambert(xLambert, yLambert float64, accuracy gpsfiles.MapAccuracy, nFileWorker int) (*Tile, int, int) {
	var t Tile

	folderPath := TILE_FOLDER_PATH + string(accuracy) + "/"

	// Getting the right file
	path, er, xLambertLL, yLambertLL := gpsfiles.ComputeTilePathFromLambert(xLambert, yLambert, folderPath, string(accuracy))

	t.XLambertLL = xLambertLL
	t.YLambertLL = yLambertLL
	if accuracy == gpsfiles.ACCURACY_1 {
		t.CellSize = 1
	} else if accuracy == gpsfiles.ACCURACY_5 {
		t.CellSize = 5
	} else if accuracy == gpsfiles.ACCURACY_25 {
		t.CellSize = 25
	}

	if er != nil {
		fmt.Println(er)
		return &t, -1, -1
	}

	fmt.Printf("Tile is: %s\n", path)

	x, y := LambertToIndices(t.XLambertLL, t.YLambertLL, xLambert, yLambert, t.CellSize)

	if x == -1 || y == -1 {
		fmt.Printf("Erreur sur le calcul de la case départ\n")
		return &t, -1, -1
	}

	fmt.Printf("x: %d, y: %d\n", x, y)
	t.CreateMatrixParallel(path, nFileWorker)

	return &t, x, y
}

func fastParseFloat(b []byte) float64 {
	var val float64
	var i int

	if len(b) == 0 {
		return 0
	}

	// Partie entière
	for ; i < len(b); i++ {
		if b[i] == '.' {
			i++
			break
		}
		val = val*10 + float64(b[i]-'0')
	}

	// Partie décimale
	div := 1.0
	for ; i < len(b); i++ {
		div *= 10
		val += float64(b[i]-'0') / div
	}

	return val
}

func (t *Tile) CreatePotentiallyReachable(d float64, startAltitude float64) {
	var potentiallyReachable [MATRIX_SIZE][MATRIX_SIZE]bool
	var reachable [MATRIX_SIZE][MATRIX_SIZE]bool

	for i := 0; i < MATRIX_SIZE; i++ {
		for j := 0; j < MATRIX_SIZE; j++ {
			if math.Abs((*t.Altitudes)[i][j]-startAltitude) < d {
				potentiallyReachable[i][j] = true
			} else {
				potentiallyReachable[i][j] = false
			}
		}
	}

	t.PotentiallyReachable = &potentiallyReachable
	t.Reachable = &reachable
}

// prend en argument le nom du fichier de BD de départ et les coordonnées x et y du point de départ en lambert et renvoie les indices i et j de la case de départ
func LambertToIndices(xLL, yLL, xLambert, yLambert, squareSize float64) (xMatrix, yMatrix int) {
	xMatrix = int((xLambert - xLL) / squareSize) // car entre chaque indice il y a 25m
	yMatrix = 1000 - int((yLambert-yLL)/squareSize)

	if xMatrix >= 1000 || yMatrix >= 1000 {
		fmt.Println("Coordonnées lambert en dehors de la case")
		return -1, -1
	}

	return
}

func (t *Tile) CreateMatrixParallel(path string, nFileWorker int) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	var altitudesMatrix [MATRIX_SIZE][MATRIX_SIZE]float64
	startOffset := 0
	wordsSkipped := 0
	for wordsSkipped < 12 && startOffset < len(data) {
		for startOffset < len(data) && (data[startOffset] <= ' ') {
			startOffset++
		}
		hasWord := false
		for startOffset < len(data) && (data[startOffset] > ' ') {
			startOffset++
			hasWord = true
		}
		if hasWord {
			wordsSkipped++
		}
	}

	payload := data[startOffset:]

	lines := bytes.Split(payload, []byte{'\n'})

	var validLines [][]byte
	for _, l := range lines {
		if len(bytes.TrimSpace(l)) > 0 {
			validLines = append(validLines, l)
		}
	}

	var wg sync.WaitGroup
	wg.Add(nFileWorker)

	rowsPerWorker := len(validLines) / nFileWorker
	if rowsPerWorker == 0 {
		rowsPerWorker = 1
	}

	for w := 0; w < nFileWorker; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == nFileWorker-1 {
			endRow = len(validLines)
		}

		go func(linesChunk [][]byte, startIdx int) {
			defer wg.Done()

			for i, line := range linesChunk {
				matrixRowIndex := startIdx + i
				if matrixRowIndex >= MATRIX_SIZE {
					break
				}

				col := 0
				pos := 0
				lineLen := len(line)

				for pos < lineLen && col < MATRIX_SIZE {
					for pos < lineLen && line[pos] <= ' ' {
						pos++
					}
					if pos >= lineLen {
						break
					}

					startNum := pos
					for pos < lineLen && line[pos] > ' ' {
						pos++
					}

					val := fastParseFloat(line[startNum:pos])

					altitudesMatrix[col][matrixRowIndex] = val
					col++
				}
			}
		}(validLines[startRow:endRow], startRow)
	}

	wg.Wait()
	t.Altitudes = &altitudesMatrix
}
