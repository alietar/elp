package tileutils

import (
	"maps"
	"math"
	"slices"
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
	cellSize := gpsfiles.ParseAccuracyFloat(accuracy)
	tileSize := 1000.0 * cellSize

	tileX := math.Floor((xLambert+cellSize/2)/tileSize)*tileSize - cellSize/2
	tileY := math.Floor((yLambert-cellSize/2)/tileSize)*tileSize + cellSize/2

	key := TileKey{X: tileX, Y: tileY}

	tc.Lock()
	// If tile already exists in cache, recover it
	if tile, exists := tc.cache[key]; exists {
		tc.Unlock()

		x, y := LambertToIndices(tile.XLambertLL, tile.YLambertLL, xLambert, yLambert, tile.CellSize)
		return tile, x, y
	}

	// If it doesn't exists, create it
	tile, x, y := NewTileFromLambert(xLambert, yLambert, accuracy)

	// Return after an error before adding to cache
	if x == -1 && y == -1 || tile.Altitudes == nil {
		tc.Unlock()
		return nil, -1, -1
	}

	tc.cache[key] = tile
	tc.Unlock()

	return tile, x, y
}

func (tc *TileCache) GetValuesSlice() []*Tile {
	return slices.Collect(maps.Values(tc.cache))
}
