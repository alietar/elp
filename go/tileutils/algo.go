package tileutils

import (
	"sync"
)

type Point struct {
	x, y int
}

func FindNeighbors(t *Tile, startX, startY int, wg *sync.WaitGroup, exploreAdj chan [2]float64) {
	defer wg.Done()

	// Pile pour l'algo itératif
	stack := make([]Point, 0, 1000)

	// ASTUCE : Puisque startX, startY est déjà marqué "Reachable=true" par le main,
	// on ne l'ajoute pas à la pile pour traitement "standard".
	// On ajoute directement ses VOISINS valides à la pile.

	// On appelle une petite fonction helper (ou on copie le bloc if) juste pour les voisins du départ
	pushNeighbors(t, startX, startY, &stack)

	// Boucle principale de Flood Fill
	for len(stack) > 0 {
		// Pop
		idx := len(stack) - 1
		p := stack[idx]
		stack = stack[:idx]

		x, y := p.x, p.y

		// 1. Check rapide (Optimisation lecture)
		// Note: PotentiallyReachable est constant une fois créé, pas besoin de lock pour lire
		if !t.PotentiallyReachable[x][y] {
			continue
		}

		// 2. Check Atomique (Ecriture)
		t.Mutex.Lock()
		if t.Reachable[x][y] {
			t.Mutex.Unlock()
			continue // Déjà fait par quelqu'un d'autre entre temps
		}
		t.Reachable[x][y] = true // Marquage
		t.Mutex.Unlock()

		// 3. Gestion des bordures (Code existant)
		if x == 0 || x == MATRIX_SIZE-1 || y == 0 || y == MATRIX_SIZE-1 {
			coord := [2]float64{
				t.XLambertLL + t.CellSize*float64(x),
				t.YLambertLL + t.CellSize*float64(1000-y),
			}

			// Ajustements précis des coordonnées (comme dans ton code original)
			if x == 0 {
				coord[0] -= t.CellSize * 1.2
			}
			if x == MATRIX_SIZE-1 {
				coord[0] += t.CellSize * 1.2
			}
			if y == 0 {
				coord[1] += t.CellSize * 1.2
			}
			if y == MATRIX_SIZE-1 {
				coord[1] -= t.CellSize * 1.2
			}

			// Envoi au coordinateur pour explorer la tuile d'à côté
			wg.Add(1)
			exploreAdj <- coord
		}

		// 4. Ajout des voisins
		pushNeighbors(t, x, y, &stack)
	}
}

func pushNeighbors(t *Tile, x, y int, stack *[]Point) {
	// On ajoute juste à la pile, on ne vérifie pas "Reachable" ici,
	// ce sera vérifié au moment du "Pop" pour garantir la thread-safety.
	if x > 0 && t.PotentiallyReachable[x-1][y] {
		*stack = append(*stack, Point{x - 1, y})
	}
	if x < MATRIX_SIZE-1 && t.PotentiallyReachable[x+1][y] {
		*stack = append(*stack, Point{x + 1, y})
	}
	if y > 0 && t.PotentiallyReachable[x][y-1] {
		*stack = append(*stack, Point{x, y - 1})
	}
	if y < MATRIX_SIZE-1 && t.PotentiallyReachable[x][y+1] {
		*stack = append(*stack, Point{x, y + 1})
	}
}
