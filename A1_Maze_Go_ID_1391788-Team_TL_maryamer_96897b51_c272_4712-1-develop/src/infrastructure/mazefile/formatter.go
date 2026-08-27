package mazefile

import (
	"errors"
	"fmt"
	"io"

	"maze/domain"
)

type Formatter struct{}

func NewFormatter() Formatter {
	return Formatter{}
}

func (Formatter) Format(writer io.Writer, maze domain.Maze) error {
	if writer == nil {
		return errors.New("приемник данных не задан")
	}

	if maze.Rows < 1 || maze.Rows > maxDimension || maze.Cols < 1 || maze.Cols > maxDimension {
		return fmt.Errorf("размер лабиринта должен быть от 1x1 до %dx%d", maxDimension, maxDimension)
	}

	if len(maze.RightWalls) != maze.Rows || len(maze.BottomWalls) != maze.Rows {
		return errors.New("некорректный размер матриц стен")
	}

	return printFileFormat(writer, maze)
}

func printFileFormat(writer io.Writer, m domain.Maze) error {
	if _, err := fmt.Fprintf(writer, "%d %d\n", m.Rows, m.Cols); err != nil {
		return err
	}

	// Матрица правых стен
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.RightWalls[r][c] {
				if _, err := fmt.Fprint(writer, "1"); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprint(writer, "0"); err != nil {
					return err
				}
			}
			if c < m.Cols-1 {
				if _, err := fmt.Fprint(writer, " "); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(writer); err != nil { // пустая строка-разделитель
		return err
	}

	// Матрица нижних стен
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.BottomWalls[r][c] {
				if _, err := fmt.Fprint(writer, "1"); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprint(writer, "0"); err != nil {
					return err
				}
			}
			if c < m.Cols-1 {
				if _, err := fmt.Fprint(writer, " "); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	return nil
}
