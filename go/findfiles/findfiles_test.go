package findfiles

import (
	"fmt"
	"testing"
)

func TestFromGpsWgs84ToLambert93(t *testing.T) {
	// Coordonnées GPS WGS84 (Paris)
	longitude := "2.3522"
	latitude := "48.8566"

	x, y, err := FromGpsWgs84ToLambert93(longitude, latitude)
	if err != nil {
		fmt.Println("Erreur de conversion :", err)
		return
	}

	fmt.Printf("Lambert-93 → X = %.2f m | Y = %.2f m\n", x, y)
}
