package algo

import (
	"fmt"
	"testing"
)

func BenchmarkTrouverVoisins(b *testing.B) {
	matriceTest := [5][5]float64{
		{341.76, 0, 0, 340.95, 0},
		{0, 341.17, 0, 340.75, 340.50},
		{0, 341.01, 340.79, 340.64, 0},
		{341.16, 340.89, 0, 0, 0},
		{341.08, 0, 340.60, 340.35, 339.97}}

	var matriceRes [5][5]float64

	for b.Loop() {
		matriceRes = TrouverVoisins(matriceTest, 1, 1)
	}

	fmt.Println(matriceRes)
}
