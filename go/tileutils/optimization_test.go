package tileutils

import (
	"fmt"
	"testing"
)

func TestOptimizeSquares(t *testing.T) {

	// --- Matrice Position de test ---
	position := [][]bool{
		{false, false, false, false, false, false},
		{false, true, true, true, false, false},
		{false, true, true, true, false, false},
		{false, true, true, true, true, true},
		{false, false, false, true, true, true},
		{false, false, false, false, false, false},
	}

	// --- Paramètres Lambert ---
	xLambertLL := 700000.0
	yLambertLL := 6600000.0
	cellSize := 25

	// --- Appel de la fonction à tester ---
	opt := OptimizeSquares(position, xLambertLL, yLambertLL, cellSize)

	// --- Vérifications basiques ---
	if len(opt.Squares) == 0 {
		t.Fatalf("Aucun carré détecté, résultat inattendu")
	}

	// --- Affichage pour vérification visuelle ---
	fmt.Println("Carrés optimisés détectés :")
	fmt.Println("--------------------------------")

	for i, sq := range opt.Squares {
		fmt.Printf("Carré #%d\n", i+1)
		fmt.Printf("  Coin haut-gauche : (%d, %d)\n", sq.X, sq.Y)
		fmt.Printf("  Taille           : %d x %d\n", sq.Size, sq.Size)
		fmt.Printf("  Centre Lambert   : (%.2f , %.2f)\n",
			sq.CenterX, sq.CenterY)
		fmt.Println()
	}
}
