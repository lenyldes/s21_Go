package eller

import (
	"math/rand"
	"testing"

	"maze/domain"
)

// TestGenerate_InvalidDimensions проверяет все граничные некорректные размеры
func TestGenerate_InvalidDimensions(t *testing.T) {
	testCases := []struct {
		name string
		rows int
		cols int
	}{
		{"zero rows", 0, 10},
		{"zero cols", 10, 0},
		{"both zero", 0, 0},
		{"negative rows", -1, 5},
		{"negative cols", 5, -1},
		{"both negative", -10, -10},
		{"rows above max", 51, 10},
		{"cols above max", 10, 51},
		{"both above max", 51, 51},
		{"huge dimensions", 1000, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Generate(tc.rows, tc.cols)
			if err != ErrInvalidDimensions {
				t.Errorf("Generate(%d, %d) expected ErrInvalidDimensions, got %v", tc.rows, tc.cols, err)
			}
		})
	}
}

// TestGenerate_BoundarySizes проверяет критические и граничные допустимые размеры
func TestGenerate_BoundarySizes(t *testing.T) {
	boundarySizes := []struct {
		name string
		rows int
		cols int
	}{
		{"min single cell 1x1", 1, 1},
		{"single row min-cols 1x2", 1, 2},
		{"single col min-rows 2x1", 2, 1},
		{"single row max-cols 1x50", 1, 50},
		{"single col max-rows 50x1", 50, 1},
		{"minimal square 2x2", 2, 2},
		{"two rows max cols 2x50", 2, 50},
		{"max rows two cols 50x2", 50, 2},
		{"medium square 10x10", 10, 10},
		{"medium rect 10x30", 10, 30},
		{"medium rect 30x10", 30, 10},
		{"max square 50x50", 50, 50},
	}

	for _, tc := range boundarySizes {
		t.Run(tc.name, func(t *testing.T) {
			maze, err := Generate(tc.rows, tc.cols)
			if err != nil {
				t.Fatalf("Generate(%d, %d) unexpected error: %v", tc.rows, tc.cols, err)
			}
			validateMatrixIntegrity(t, maze, tc.rows, tc.cols)
			validatePerfectMaze(t, maze)
		})
	}
}

// TestGenerate_FuzzStress 500 рандомных генераций
func TestGenerate_FuzzStress(t *testing.T) {
	iterations := 500
	for i := 0; i < iterations; i++ {
		rows := rand.Intn(50) + 1
		cols := rand.Intn(50) + 1

		maze, err := Generate(rows, cols)
		if err != nil {
			t.Fatalf("Iteration %d: Generate(%d, %d) failed: %v", i, rows, cols, err)
		}
		validateMatrixIntegrity(t, maze, rows, cols)
		validatePerfectMaze(t, maze)
	}
}

// TestGenerate_Randomness проверяем что генерация не проходит однотипно
func TestGenerate_Randomness(t *testing.T) {
	rows, cols := 15, 15
	m1, err1 := Generate(rows, cols)
	m2, err2 := Generate(rows, cols)
	if err1 != nil || err2 != nil {
		t.Fatalf("Generation failed: %v, %v", err1, err2)
	}

	// Сравниваем, что хотя бы одна внутренняя стенка отличается
	different := false
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if m1.RightWalls[r][c] != m2.RightWalls[r][c] || m1.BottomWalls[r][c] != m2.BottomWalls[r][c] {
				different = true
				break
			}
		}
		if different {
			break
		}
	}

	if !different {
		t.Errorf("Generate(%d, %d) produced identical mazes twice in a row", rows, cols)
	}
}

// validateMatrixIntegrity проверяет размеры выделенной памяти и срезов
func validateMatrixIntegrity(t *testing.T, m domain.Maze, expectedRows, expectedCols int) {
	t.Helper()
	if m.Rows != expectedRows || m.Cols != expectedCols {
		t.Fatalf("Maze dimensions mismatch: got %dx%d, expected %dx%d", m.Rows, m.Cols, expectedRows, expectedCols)
	}
	if len(m.RightWalls) != expectedRows || len(m.BottomWalls) != expectedRows {
		t.Fatalf("Matrix rows slice length mismatch: right=%d, bottom=%d, expected=%d",
			len(m.RightWalls), len(m.BottomWalls), expectedRows)
	}
	for r := 0; r < expectedRows; r++ {
		if len(m.RightWalls[r]) != expectedCols || len(m.BottomWalls[r]) != expectedCols {
			t.Fatalf("Row %d cols slice length mismatch: right=%d, bottom=%d, expected=%d",
				r, len(m.RightWalls[r]), len(m.BottomWalls[r]), expectedCols)
		}
	}
}

// validatePerfectMaze проверяет строгие свойства идеального лабиринта:
// 1. Границы лабиринта закрыты справа и снизу.
// 2. Количество рёбер (проходов) равно количеству вершин - 1 (E = V - 1).
// 3. Лабиринт связный (все ячейки достижимы из (0,0)).
// 4. В графе нет ни одного цикла/петли.
func validatePerfectMaze(t *testing.T, m domain.Maze) {
	t.Helper()

	rows, cols := m.Rows, m.Cols
	totalCells := rows * cols

	// 1. Проверяем внешние границы
	for r := 0; r < rows; r++ {
		if !m.RightWalls[r][cols-1] {
			t.Errorf("Right boundary wall missing at row %d", r)
		}
	}
	for c := 0; c < cols; c++ {
		if !m.BottomWalls[rows-1][c] {
			t.Errorf("Bottom boundary wall missing at col %d", c)
		}
	}

	// Строим граф лабиринта
	adj := make(map[int][]int)
	edgeCount := 0

	toID := func(r, c int) int {
		return r*cols + c
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			u := toID(r, c)

			// Проход вправо
			if c < cols-1 && !m.RightWalls[r][c] {
				v := toID(r, c+1)
				adj[u] = append(adj[u], v)
				adj[v] = append(adj[v], u)
				edgeCount++
			}

			// Проход вниз
			if r < rows-1 && !m.BottomWalls[r][c] {
				v := toID(r+1, c)
				adj[u] = append(adj[u], v)
				adj[v] = append(adj[v], u)
				edgeCount++
			}
		}
	}

	// 2. Проверка свойства остовного дерева: E = V - 1
	expectedEdges := totalCells - 1
	if edgeCount != expectedEdges {
		t.Errorf("Maze %dx%d has %d passages (edges), expected %d (V-1). Loops or isolations exist!",
			rows, cols, edgeCount, expectedEdges)
	}

	// 3 и 4. Проверка связности и отсутствия циклов через BFS
	visited := make(map[int]bool)
	parent := make(map[int]int)

	queue := []int{0}
	visited[0] = true
	parent[0] = -1

	hasCycle := false

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, neighbor := range adj[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = curr
				queue = append(queue, neighbor)
			} else if parent[curr] != neighbor {
				hasCycle = true
			}
		}
	}

	if hasCycle {
		t.Errorf("Maze %dx%d contains a loop / cycle!", rows, cols)
	}

	if len(visited) != totalCells {
		t.Errorf("Maze %dx%d is not fully connected: visited %d/%d cells (isolated areas exist!)",
			rows, cols, len(visited), totalCells)
	}
}
