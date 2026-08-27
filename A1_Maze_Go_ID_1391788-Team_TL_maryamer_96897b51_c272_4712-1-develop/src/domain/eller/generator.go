package eller

import (
	"errors"
	"math/rand"

	"maze/domain"
)

var (
	ErrInvalidDimensions = errors.New("invalid maze dimensions: rows and cols must be between 1 and 50")
)

const (
	minDimension = 1
	maxDimension = 50
)

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

func (Generator) Generate(rows, cols int) (domain.Maze, error) {
	return Generate(rows, cols)
}

func Generate(rows, cols int) (domain.Maze, error) {
	if rows < minDimension || rows > maxDimension || cols < minDimension || cols > maxDimension {
		return domain.Maze{}, ErrInvalidDimensions
	}

	rightWalls := make([][]bool, rows)
	bottomWalls := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		rightWalls[i] = make([]bool, cols)
		bottomWalls[i] = make([]bool, cols)
	}

	sets := make([]int, cols)
	nextSetID := 1

	for r := 0; r < rows-1; r++ {
		nextSetID = assignSets(sets, nextSetID)

		// Расстановка правых стенок
		rightWalls[r][cols-1] = true
		for c := 0; c < cols-1; c++ {
			if sets[c] == sets[c+1] {
				rightWalls[r][c] = true
				continue
			}

			if randWall() {
				rightWalls[r][c] = true
			} else {
				rightWalls[r][c] = false
				mergeSets(sets, sets[c+1], sets[c])
			}
		}

		// Расстановка нижних стенок
		// Сначала случайно расставляем нижние стенки всем ячейкам
		for c := 0; c < cols; c++ {
			bottomWalls[r][c] = randWall()
		}

		// Потом проверяем что у каждого множества есть хотя бы один проход вниз
		// Если у какго то множества все нижние стенки закрыты то просто открываем случайную
		setCells := make(map[int][]int) // k - id множества, v - индексы ячеек
		for c := 0; c < cols; c++ {
			setCells[sets[c]] = append(setCells[sets[c]], c)
		}

		for _, cells := range setCells {
			hasDoor := false
			for _, c := range cells {
				if !bottomWalls[r][c] {
					hasDoor = true
					break
				}
			}
			if !hasDoor {
				openIndex := cells[rand.Intn(len(cells))]
				bottomWalls[r][openIndex] = false
			}
		}

		// Подготовка к следующей строке
		// Ячейки с нижней стенкой теряют принадлежность к множеству (становятся 0)
		for c := 0; c < cols; c++ {
			if bottomWalls[r][c] {
				sets[c] = 0
			}
		}
	}

	// Обработка последней строки
	lastRow := rows - 1
	nextSetID = assignSets(sets, nextSetID)

	// Всем ячейкам последней строки ставим нижнюю стенку (граница лабиринта)
	for c := 0; c < cols; c++ {
		bottomWalls[lastRow][c] = true
	}
	rightWalls[lastRow][cols-1] = true // Крайняя правая граница всегда закрыта

	// Устраняем изолированные области: убираем правые стенки между РАЗНЫМИ множествами
	for c := 0; c < cols-1; c++ {
		if sets[c] != sets[c+1] {
			rightWalls[lastRow][c] = false
			mergeSets(sets, sets[c+1], sets[c])
		} else {
			rightWalls[lastRow][c] = true
		}
	}

	return domain.NewMaze(rows, cols, rightWalls, bottomWalls), nil
}

func assignSets(sets []int, nextSetID int) int {
	for c := range sets {
		if sets[c] == 0 {
			sets[c] = nextSetID
			nextSetID++
		}
	}
	return nextSetID
}

func randWall() bool {
	return rand.Intn(2) == 1
}

func mergeSets(sets []int, targetSet, sourceSet int) {
	for i := range sets {
		if sets[i] == targetSet {
			sets[i] = sourceSet
		}
	}
}
