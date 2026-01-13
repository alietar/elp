package algo

import (
	lgo "github.com/yageek/lambertgo"
)

// crée une structure Coordonnee pour stocker les coordonnées GPS
type Coordonnee struct {
	Lat  float64 `json:"lat"` // gère l'export en JSON
	Lng  float64 `json:"lng"`
	Size int     `json:"size"`
}

// Convertit des coordonnées Lambert93 en coordonnées GPS WGS84
func FromLambert93ToGpsWgs84(x, y float64) (float64, float64, error) {

	// Crée le point avec les coordonnées Lambert93
	point := &lgo.Point{
		X:    x,
		Y:    y,
		Z:    0,
		Unit: lgo.Meter, // Préciser que c'est en mètres
	}

	point.ToWGS84(lgo.Lambert93) // Conversion Lambert-93 → WGS84

	// Convertion de résultat en Radians vers des Degrés pour Leaflet de JS
	point.ToDegree()

	// point.Y est la Latitude, point.X est la Longitude
	return point.Y, point.X, nil
}
