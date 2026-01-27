package tileutils

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/alietar/elp/go/gpsfiles"
)

func NewTileFromLambert(xLambert, yLambert float64, accuracy gpsfiles.MapAccuracy) (*Tile, int, int) {
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
	// t.CreateMatrix(path)
	// t.CreateMatrixOptimized(path)
	// t.CreateMatrixHighPerf(path)
	t.CreateMatrixParallel(path)

	return &t, x, y
}

// prend en argument le nom du fichier de BD de départ et renvoie une matrice 1000x1000 des altitudes
func (t *Tile) CreateMatrix(path string) {
	data, err := os.ReadFile(path) // lire le fichier en question, data est en byte

	if err != nil {
		fmt.Println("Erreur lecture base de données:", err)
		return
	}

	var altitudesMatrix [MATRIX_SIZE][MATRIX_SIZE]float64

	altitudes := strings.Fields(string(data))[12:] // on sépare chaque element du fichier (sépare dès qu'il y a un espace ou retour à la ligne)
	// les altitudes démarrent au 13e élément

	k := 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			altitude, err := strconv.ParseFloat(altitudes[k], 64) // conversion de string en float64
			if err != nil {
				fmt.Println("Erreur lors de la récupération de l'altitude", err)
			}
			altitudesMatrix[j][i] = altitude
			k += 1
		}
	}

	t.Altitudes = &altitudesMatrix
}

func (t *Tile) CreateMatrixOptimized(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Erreur ouverture fichier:", err)
		return
	}
	defer file.Close()

	// Scanner pour lire mot par mot (espaces et sauts de ligne sont gérés automatiquement)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	// Sauter les 12 premiers éléments (méta-données)
	for k := 0; k < 12; k++ {
		if !scanner.Scan() {
			return // Fin de fichier prématurée
		}
	}

	var altitudesMatrix [MATRIX_SIZE][MATRIX_SIZE]float64

	for i := 0; i < MATRIX_SIZE; i++ {
		for j := 0; j < MATRIX_SIZE; j++ {
			if scanner.Scan() {
				// scanner.Text() alloue une petite string, mais c'est vite libéré
				val := fastParseFloat(scanner.Bytes())
				// ATTENTION : J'ai inversé [j][i] en [i][j] pour l'exemple (voir point 3 ci-dessous)
				altitudesMatrix[j][i] = val
			}
		}
	}

	t.Altitudes = &altitudesMatrix
}

func (t *Tile) CreateMatrixHighPerf(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var altitudesMatrix [MATRIX_SIZE][MATRIX_SIZE]float64

	// Position actuelle dans le buffer 'data'
	pos := 0

	// Fonction helper inline pour avancer jusqu'au prochain mot
	// Retourne le début et la fin du mot
	nextWord := func() (int, int) {
		// Ignorer les espaces/newlines
		for pos < len(data) && (data[pos] == ' ' || data[pos] == '\n' || data[pos] == '\r') {
			pos++
		}
		start := pos
		// Trouver la fin du mot
		for pos < len(data) && data[pos] != ' ' && data[pos] != '\n' && data[pos] != '\r' {
			pos++
		}
		return start, pos
	}

	// Sauter les 12 premiers
	for k := 0; k < 12; k++ {
		nextWord()
	}

	for i := 0; i < MATRIX_SIZE; i++ {
		for j := 0; j < MATRIX_SIZE; j++ {
			start, end := nextWord()
			if start >= end {
				break
			}

			// Astuce Zero-Copy : Créer une string depuis le slice de bytes sans allocation
			// Nécessite Go 1.20+ pour unsafe.String, sinon utiliser l'ancien trick *reflect.StringHeader
			valStr := unsafe.String(&data[start], end-start)

			val, _ := strconv.ParseFloat(valStr, 64)
			altitudesMatrix[j][i] = val
		}
	}

	t.Altitudes = &altitudesMatrix
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

func (t *Tile) CreateMatrixParallel(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	var altitudesMatrix [MATRIX_SIZE][MATRIX_SIZE]float64

	// 1. Sauter le header (les 12 premiers mots)
	// On cherche simplement le début des données brutes
	startOffset := 0
	wordsSkipped := 0
	for wordsSkipped < 12 && startOffset < len(data) {
		// Avance tant qu'il y a des espaces
		for startOffset < len(data) && (data[startOffset] <= ' ') {
			startOffset++
		}
		// Avance tant qu'il y a du texte (le mot)
		hasWord := false
		for startOffset < len(data) && (data[startOffset] > ' ') {
			startOffset++
			hasWord = true
		}
		if hasWord {
			wordsSkipped++
		}
	}

	// On ne garde que la partie utile du fichier
	payload := data[startOffset:]

	// 2. Préparer le découpage par lignes
	// On trouve tous les indices de retour à la ligne pour distribuer le travail
	// Astuce : On suppose que le format est stable (1000 lignes).
	// Si le fichier est un "long texte" sans lignes claires, l'approche change.
	lines := bytes.Split(payload, []byte{'\n'})

	// On filtre les lignes vides (souvent une à la fin)
	var validLines [][]byte
	for _, l := range lines {
		if len(bytes.TrimSpace(l)) > 0 {
			validLines = append(validLines, l)
		}
	}

	numWorkers := 4 // Utilise tous les cœurs disponibles
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	rowsPerWorker := len(validLines) / numWorkers
	if rowsPerWorker == 0 {
		rowsPerWorker = 1
	}

	// 3. Lancer les workers
	for w := 0; w < numWorkers; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == numWorkers-1 {
			endRow = len(validLines) // Le dernier prend le reste
		}

		go func(linesChunk [][]byte, startIdx int) {
			defer wg.Done()

			for i, line := range linesChunk {
				matrixRowIndex := startIdx + i
				if matrixRowIndex >= MATRIX_SIZE {
					break
				}

				// Parsing manuel de la ligne sans allocation string
				col := 0
				pos := 0
				lineLen := len(line)

				for pos < lineLen && col < MATRIX_SIZE {
					// Sauter espaces initiaux
					for pos < lineLen && line[pos] <= ' ' {
						pos++
					}
					if pos >= lineLen {
						break
					}

					// Début du nombre
					startNum := pos
					// Trouver fin du nombre
					for pos < lineLen && line[pos] > ' ' {
						pos++
					}

					// Conversion haute performance
					// Option A (Sûr) : val, _ := strconv.ParseFloat(unsafe.String(&line[startNum], pos-startNum), 64)
					// Option B (Risqué mais Extrêmement rapide) :
					val := fastParseFloat(line[startNum:pos])

					// Note: Inversion [col][row] vs [row][col] selon ta structure
					altitudesMatrix[col][matrixRowIndex] = val
					col++
				}
			}
		}(validLines[startRow:endRow], startRow)
	}

	wg.Wait()
	t.Altitudes = &altitudesMatrix
}
