package findfiles

import (
	lgo "github.com/yageek/lambertgo"

	"strconv"
)

func FromGpsWgs84ToLambert93(longStr string, latStr string) (float64, float64, error) {
	longitude, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return 0, 0, err
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, err
	}

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
