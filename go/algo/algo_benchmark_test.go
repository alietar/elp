package algo

import (
	"testing"
)

func BenchmarkFindingNeighbors(b *testing.B) {
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

	for b.Loop() {
		testMatrix.FindNeighbors(1, 1)
	}
}
