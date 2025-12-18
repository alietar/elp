package algo

import (
	"fmt"
	"testing"
)

func TestTrouverVoisins(t *testing.T) {
	matrice := [5][5]float64{
		{341.76, 0, 0, 340.95, 0},
		{0, 341.17, 0, 340.75, 340.50},
		{0, 341.01, 340.79, 340.64, 0},
		{341.16, 340.89, 0, 0, 0},
		{341.08, 0, 340.60, 340.35, 339.97}}

	afficherMatrice(matrice)

	fmt.Printf("\n---\n\n")

	afficherMatrice(TrouverVoisins(matrice, 1, 1))
	/*t.Run341.08 340.82 340.60 340.35 339.97("saying hello to people", func(t *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris"
	})*/
}
