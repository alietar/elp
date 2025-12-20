package algo

import (
	"fmt"
	"testing"
)

func TestTrouverVoisins(t *testing.T) {
	matriceTest := [5][5]float64{
		{341.76, 0, 0, 340.95, 0},
		{0, 341.17, 0, 340.75, 340.50},
		{0, 341.01, 340.79, 340.64, 0},
		{341.16, 340.89, 0, 0, 0},
		{341.08, 0, 340.60, 340.35, 339.97}}

	// afficherMatrice(matrice)

	// fmt.Printf("\n---\n\n")

	// afficherMatrice(TrouverVoisins(matrice, 1, 1))

	t.Run("matrice_taille_5_situation_1", func(t *testing.T) {
		got := TrouverVoisins(matriceTest, 0, 0)
		want := [5][5]float64{
			{341.76, 0, 0, 340.95, 0},
			{0, 341.17, 0, 340.75, 340.50},
			{0, 341.01, 340.79, 340.64, 0},
			{341.16, 340.89, 0, 0, 0},
			{341.08, 0, 340.60, 340.35, 339.97}}

		if got != want {
			fmt.Println("Got :")
			afficherMatrice(got)
			fmt.Println("Want :")
			afficherMatrice(want)
			t.Errorf("aaa")
		}
	})
}

func evaluerReponse(t testing.TB, got, want [5][5]float64) {
	t.Helper()

	if got != want {
		fmt.Println("Got :")
		afficherMatrice(got)
		fmt.Println("Want :")
		afficherMatrice(want)
		t.Errorf("La matrice réponse n'est pas la bonne")
	}
}

/*
func TestCreation_matrice(t *testing.T) {
	// Creation_matrice()
}
*/
