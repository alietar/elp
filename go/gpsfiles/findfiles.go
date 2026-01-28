package gpsfiles

import (
	"math"

	lgo "github.com/yageek/lambertgo"

	"fmt"
)

func ConvertWgs84ToLambert93(longitude, latitude float64) (float64, float64, error) {
	// Point WGS84 en degrés
	point := &lgo.Point{
		X:    longitude,
		Y:    latitude,
		Z:    0,
		Unit: lgo.Degree,
	}

	// degrés → radians
	point.ToRadian()

	// WGS84 → Lambert-93 (EPSG:2154)
	point.ToLambert(lgo.Lambert93)

	return point.X, point.Y, nil
}

func ConvertLambert93ToWgs84(x, y float64) (float64, float64, error) {
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

func ComputeTilePathFromLambert(x, y float64, path string, accuracy string) (string, error, float64, float64) {
	cellSize := ParseAccuracyFloat(MapAccuracy(accuracy))

	fileX := math.Floor((x+cellSize/2)/(cellSize*1000)) * cellSize
	fileY := math.Floor((y-cellSize/2)/(cellSize*1000))*cellSize + cellSize

	path = fmt.Sprintf("%s%s_%04.0f_%04.0f.asc", path, string(accuracy), fileX, fileY)

	// Offsetting to get the right lower left coordinates
	xll := fileX*1000 - cellSize/2
	yll := (fileY-cellSize)*1000 + cellSize/2

	// Bounds math
	xmin := xll
	xmax := xll + 1000*cellSize
	ymin := yll
	ymax := yll + 1000*cellSize

	// Double checking if coordinates are in tile
	if x >= xmin && x <= xmax && y >= ymin && y <= ymax {
		return path, nil, xll, yll
	}

	return "", fmt.Errorf("coordinate not found in any file"), -1, -1
}
