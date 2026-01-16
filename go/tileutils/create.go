package tileutils

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/alietar/elp/go/gpsfiles"
)

func NewTileFromLambert(xLambert, yLambert float64) (*Tile, int, int) {
	var t Tile

	// Getting the right file
	path, er, xLambertLL, yLambertLL := gpsfiles.GetFileForMyCoordinate(xLambert, yLambert, TILE_FOLDER_PATH)

	t.XLambertLL = xLambertLL
	t.YLambertLL = yLambertLL

	if er != nil {
		fmt.Println(er)
		return &t, -1, -1
	}

	path = TILE_FOLDER_PATH + path

	fmt.Printf("Tile is: %s\n", path)

	x, y := LambertToIndices(t.XLambertLL, t.YLambertLL, xLambert, yLambert)

	if x == -1 || y == -1 {
		fmt.Printf("Erreur sur le calcul de la case départ\n")
		return &t, -1, -1
	}

	fmt.Printf("x: %d, y: %d\n", x, y)
	t.CreateMatrix(path)

	// startAltitude := (*t.Altitudes)[x][y]
	// startAltitude := 169.75

	// t.CreatePotentiallyReachable(5, startAltitude)

	// fmt.Printf("Start Altitude : %f\n", startAltitude)

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
func LambertToIndices(xLL, yLL, xLambert, yLambert float64) (xMatrix, yMatrix int) {
	xMatrix = int((xLambert - xLL) / 25) // car entre chaque indice il y a 25m
	yMatrix = 1000 - int((yLambert-yLL)/25)

	if xMatrix >= 1000 || yMatrix >= 1000 {
		fmt.Println("Coordonnées lambert en dehors de la case")
		return -1, -1
	}

	return
}
