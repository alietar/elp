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

	t.Run("matrice_taille_5_situation_1", func(t *testing.T) {
		got := TrouverVoisins(matriceTest, 0, 0)
		want := [5][5]float64{
			{341.76, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0}}

		evaluerReponse(t, got, want)
	})

	t.Run("matrice_taille_5_situation_2", func(t *testing.T) {
		got := TrouverVoisins(matriceTest, 1, 1)
		want := [5][5]float64{
			{0, 0, 0, 340.95, 0},
			{0, 341.17, 0, 340.75, 340.50},
			{0, 341.01, 340.79, 340.64, 0},
			{341.16, 340.89, 0, 0, 0},
			{341.08, 0, 0, 0, 0}}

		evaluerReponse(t, got, want)
	})

	t.Run("matrice_taille_5_situation_3", func(t *testing.T) {
		got := TrouverVoisins(matriceTest, 4, 3)
		want := [5][5]float64{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 340.60, 340.35, 339.97}}

		evaluerReponse(t, got, want)
	})

	t.Run("matrice_taille_5_situation_4", func(t *testing.T) {
		got := TrouverVoisins(matriceTest, 1, 0)
		want := [5][5]float64{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0}}

		evaluerReponse(t, got, want)
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

func TestCreationMatrice(t *testing.T) {
	var matriceVide [1000][1000]float64

	matriceResultat := CreationMatrice("../bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/BDALTIV2_25M_FXX_0775_6550_MNT_LAMB93_IGN69.asc")
	if matriceResultat == matriceVide {
		t.Errorf("Erreur à la création de matrice")
	}
}

func TestPointsAtteignables(t *testing.T) {
	matrice := CreationMatrice("../bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/BDALTIV2_25M_FXX_0775_6550_MNT_LAMB93_IGN69.asc")
	PointsAtteignables(10, 13, 14, matrice)
}

func TestCaseDepart(t *testing.T) {
	x, y := CaseDepart(10, 10, "../bd/1_DONNEES_LIVRAISON_2024-02-00018/BDALTIV2_MNT_25M_ASC_LAMB93_IGN69_D069/BDALTIV2_25M_FXX_0775_6550_MNT_LAMB93_IGN69.asc")

	if x == -1 || y == -1 {
		t.Errorf("Erreur pour trouver la case départ")
	}

	fmt.Printf("x: %d\ny: %d\n", x, y)
}
