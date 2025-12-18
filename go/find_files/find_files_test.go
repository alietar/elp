package find_files

import (
	"testing"
	"fmt"
)

func Testfrom_GPS_WGS84_to_Lambert_93(t *testing.T) {
	
	// Coordonnées GPS WGS84 (Paris)
	longitude := "2.3522"
	latitude := "48.8566"

	x, y, err := from_GPS_WGS84_to_Lambert_93(longitude, latitude)
	if err != nil {
		fmt.Println("Erreur de conversion :", err)
		return
	}

	fmt.Printf("Lambert-93 → X = %.2f m | Y = %.2f m\n", x, y)
}
