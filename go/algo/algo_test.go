package algo

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"testing"
)

func TestFindNeighbors(t *testing.T) {
	testMatrix := Matrix{
		Size: 5,
		Data: [][]float64{
			{341.76, 0, 0, 340.95, 0},
			{0, 341.17, 0, 340.75, 340.50},
			{0, 341.01, 340.79, 340.64, 0},
			{341.16, 340.89, 0, 0, 0},
			{341.08, 0, 340.60, 340.35, 339.97},
		},
	}

	t.Run("matrice_taille_5_situation_1", func(t *testing.T) {
		got := testMatrix.FindNeighbors(0, 0)
		want := Matrix{
			Size: 5,
			Data: [][]float64{
				{341.76, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		}

		evaluateMatrices(t, got, want)
	})

	t.Run("matrice_taille_5_situation_2", func(t *testing.T) {
		got := testMatrix.FindNeighbors(1, 1)
		want := Matrix{
			Size: 5,
			Data: [][]float64{
				{0, 0, 0, 340.95, 0},
				{0, 341.17, 0, 340.75, 340.50},
				{0, 341.01, 340.79, 340.64, 0},
				{341.16, 340.89, 0, 0, 0},
				{341.08, 0, 0, 0, 0},
			},
		}

		evaluateMatrices(t, got, want)
	})

	t.Run("matrice_taille_5_situation_3", func(t *testing.T) {
		got := testMatrix.FindNeighbors(4, 3)
		want := Matrix{
			Size: 5,
			Data: [][]float64{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 340.60, 340.35, 339.97},
			},
		}

		evaluateMatrices(t, got, want)
	})

	t.Run("matrice_taille_5_situation_4", func(t *testing.T) {
		got := testMatrix.FindNeighbors(1, 0)
		want := Matrix{
			Size: 5,
			Data: [][]float64{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		}

		evaluateMatrices(t, got, want)
	})

	t.Run("matrice_aleatoire", func(t *testing.T) {
		size := 20
		randomMatrix := NewMatrix(size)

		for i := range size {
			for j := range size {
				if rand.IntN(2) == 0 {
					randomMatrix.Data[i][j] = 10
				}
			}
		}
		randomMatrix.Show()

		got := randomMatrix.FindNeighbors(0, 0)
		got.Show()
	})
}

func evaluateMatrices(t testing.TB, got, want Matrix) {
	t.Helper()

	if !reflect.DeepEqual(got.Data, want.Data) {
		fmt.Println("Got :")
		got.Show()
		fmt.Println("Want :")
		want.Show()
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

	// fmt.Printf("x: %d\ny: %d\n", x, y)
}

func TestGetDepartement(t *testing.T) {
	fmt.Println(GetDepartement(45.767712, 4.98775))
	fmt.Println(GetDepartement(47.433331, -2.08333))
}

func TestDownloadUnzipDB(t *testing.T) {
	downloadUnzipDB(69)
}
