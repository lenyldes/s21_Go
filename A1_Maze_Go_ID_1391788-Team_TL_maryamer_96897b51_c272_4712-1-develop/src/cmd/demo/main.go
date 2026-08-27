package main

import (
	"fmt"
	"os"
	"strconv"

	"maze/domain"
	"maze/domain/eller"
)

func main() {
	rows := 10
	cols := 10

	if len(os.Args) >= 3 {
		if r, err := strconv.Atoi(os.Args[1]); err == nil {
			rows = r
		}
		if c, err := strconv.Atoi(os.Args[2]); err == nil {
			cols = c
		}
	} else if len(os.Args) == 2 {
		if r, err := strconv.Atoi(os.Args[1]); err == nil {
			rows = r
			cols = r
		}
	}

	fmt.Printf("=== Генерация лабиринта %dx%d (Алгоритм Эллера) ===\n\n", rows, cols)

	maze, err := eller.Generate(rows, cols)
	if err != nil {
		fmt.Printf("Ошибка генерации: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("1. Визуализация лабиринта в терминале:")
	printAsciiMaze(maze)

	fmt.Println("\n2. Формат файла задания (README):")
	printFileFormat(maze)
}

func printAsciiMaze(m domain.Maze) {
	// Верхняя граница
	fmt.Print("+")
	for c := 0; c < m.Cols; c++ {
		fmt.Print("---+")
	}
	fmt.Println()

	for r := 0; r < m.Rows; r++ {
		// Тело ячеек и правые стены
		fmt.Print("|")
		for c := 0; c < m.Cols; c++ {
			fmt.Print("   ")
			if m.RightWalls[r][c] {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()

		// Нижние стены
		fmt.Print("+")
		for c := 0; c < m.Cols; c++ {
			if m.BottomWalls[r][c] {
				fmt.Print("---")
			} else {
				fmt.Print("   ")
			}
			fmt.Print("+")
		}
		fmt.Println()
	}
}

func printFileFormat(m domain.Maze) {
	fmt.Printf("%d %d\n", m.Rows, m.Cols)
	// Матрица правых стен
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.RightWalls[r][c] {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
			if c < m.Cols-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println() // пустая строка-разделитель

	// Матрица нижних стен
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.BottomWalls[r][c] {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
			if c < m.Cols-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
