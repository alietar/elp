package tileutils

import (
	"math"
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

// Nouvelle clé de cache : beaucoup plus efficace qu'une string
type TileKey struct {
	X float64
	Y float64
}

type TileCache struct {
	sync.Mutex
	cache map[TileKey]*Tile
}

func NewTileCache() *TileCache {
	return &TileCache{
		cache: make(map[TileKey]*Tile),
	}
}

func (tc *TileCache) GetOrLoad(xLambert, yLambert float64, accuracy gpsfiles.MapAccuracy) (*Tile, int, int) {
	// 1. DÉTERMINER LA TAILLE RÉELLE DE LA TUILE (en mètres)
	var cellSize float64
	switch accuracy {
	case gpsfiles.ACCURACY_1:
		cellSize = 1
	case gpsfiles.ACCURACY_5:
		cellSize = 5
	case gpsfiles.ACCURACY_25:
		cellSize = 25
	}

	tileSize := 1000.0 * cellSize // Ex: 25000m pour accuracy 25m

	// 2. CALCUL MATHÉMATIQUE DE L'ORIGINE (Sans toucher au disque !)
	// On arrondit à l'entier inférieur multiple de la taille de tuile
	tileX := math.Floor((xLambert+cellSize/2)/tileSize)*tileSize - cellSize/2
	tileY := math.Floor((yLambert-cellSize/2)/tileSize)*tileSize + cellSize/2

	key := TileKey{X: tileX, Y: tileY}

	// 3. CHECK CACHE (Opération purement RAM)
	tc.Lock()
	if tile, exists := tc.cache[key]; exists {
		tc.Unlock()
		// Calcul des indices locaux
		x, y := LambertToIndices(tile.XLambertLL, tile.YLambertLL, xLambert, yLambert, tile.CellSize)
		return tile, x, y
	}

	// Chargement réel
	tile, x, y := NewTileFromLambert(xLambert, yLambert, accuracy)

	if x != -1 && y != -1 {
		// On stocke avec notre clé calculée
		tc.cache[key] = tile
	}
	tc.Unlock()

	return tile, x, y
}
