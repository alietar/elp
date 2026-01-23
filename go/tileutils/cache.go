package tileutils

import (
	"sync"

	"github.com/alietar/elp/go/gpsfiles"
)

type TileCache struct {
	sync.Mutex
	cache map[string]*Tile // String for the path
}

func NewTileCache() *TileCache {
	return &TileCache{
		cache: make(map[string]*Tile),
	}
}
func (tc *TileCache) GetOrLoad(xLambert, yLambert float64, accuracy gpsfiles.MapAccuracy) (*Tile, int, int, string) {
	// On détermine le chemin du fichier pour savoir si on l'a déjà
	folderPath := "./db/" + string(accuracy) + "/"
	path, _, _, _ := gpsfiles.GetFileForMyCoordinate(xLambert, yLambert, folderPath)

	tc.Lock()
	defer tc.Unlock()

	// Si la tuile existe déjà dans le cache, on la renvoie
	if tile, exists := tc.cache[path]; exists {
		// On recalcule juste les indices locaux car ils dépendent du point d'entrée
		x, y := LambertToIndices(tile.XLambertLL, tile.YLambertLL, xLambert, yLambert, tile.CellSize)
		return tile, x, y, path
	}

	// Sinon, on crée la tuile (cette fonction doit être adaptée pour ne pas retourner un nouveau pointeur si on veut le faire proprement,
	// mais ici on utilise ta fonction existante et on stocke le résultat)
	tile, x, y := NewTileFromLambert(xLambert, yLambert, accuracy)

	// Si chargement réussi, on stocke dans le cache
	if x != -1 && y != -1 {
		tc.cache[path] = tile
	}

	return tile, x, y, path
}
