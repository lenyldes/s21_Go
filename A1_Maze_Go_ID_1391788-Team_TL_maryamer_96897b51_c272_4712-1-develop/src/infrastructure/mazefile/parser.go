package mazefile

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"maze/domain"
)

const maxDimension = 50

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

func (Parser) Parse(reader io.Reader) (domain.Maze, error) {
	if reader == nil {
		return domain.Maze{}, errors.New("источник файла не задан")
	}

	scanner := bufio.NewReader(reader)
	var rows, cols int
	if _, err := fmt.Fscan(scanner, &rows, &cols); err != nil {
		return domain.Maze{}, errors.New("не удалось прочитать размеры лабиринта")
	}
	if rows < 1 || rows > maxDimension || cols < 1 || cols > maxDimension {
		return domain.Maze{}, fmt.Errorf("размер лабиринта должен быть от 1x1 до %dx%d", maxDimension, maxDimension)
	}

	rightWalls, err := readMatrix(scanner, rows, cols, "вертикальных")
	if err != nil {
		return domain.Maze{}, err
	}
	bottomWalls, err := readMatrix(scanner, rows, cols, "горизонтальных")
	if err != nil {
		return domain.Maze{}, err
	}

	var extra string
	if _, err = fmt.Fscan(scanner, &extra); err != io.EOF {
		if err != nil {
			return domain.Maze{}, fmt.Errorf("ошибка чтения файла: %w", err)
		}
		return domain.Maze{}, errors.New("после матриц обнаружены лишние данные")
	}
	if err = validateOuterWalls(rightWalls, bottomWalls, rows, cols); err != nil {
		return domain.Maze{}, err
	}

	return domain.NewMaze(rows, cols, rightWalls, bottomWalls), nil
}

func readMatrix(reader io.Reader, rows, cols int, name string) ([][]bool, error) {
	matrix := make([][]bool, rows)
	for row := 0; row < rows; row++ {
		matrix[row] = make([]bool, cols)
		for col := 0; col < cols; col++ {
			var value int
			if _, err := fmt.Fscan(reader, &value); err != nil {
				return nil, fmt.Errorf("матрица %s стен: ожидалось %d значений", name, rows*cols)
			}
			if value != 0 && value != 1 {
				return nil, fmt.Errorf("матрица %s стен: значение [%d,%d] должно быть 0 или 1", name, row+1, col+1)
			}
			matrix[row][col] = value == 1
		}
	}
	return matrix, nil
}

func validateOuterWalls(right, bottom [][]bool, rows, cols int) error {
	for row := 0; row < rows; row++ {
		if !right[row][cols-1] {
			return fmt.Errorf("лабиринт не замкнут справа: строка %d", row+1)
		}
	}
	for col := 0; col < cols; col++ {
		if !bottom[rows-1][col] {
			return fmt.Errorf("лабиринт не замкнут снизу: столбец %d", col+1)
		}
	}
	return nil
}
