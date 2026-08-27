package application

import (
	"errors"
	"maze/domain"
)

const visited = 1
const deadend = 2
const unvisited = 0

func Solve(maze domain.Maze, start, end domain.Coords) ([][]bool, error) {
	var path = make([][]int, maze.Rows)

	for i := range path {
		path[i] = make([]int, maze.Cols)
	}

	path, ok, err := findPath(path, maze, start, end)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("impossible to find the way")
	}

	return toBooleanPath(path), nil
}

// Базовые случаи: выход, тупик, уже посещенная клетка
func findPath(path [][]int, maze domain.Maze, cell, end domain.Coords) ([][]int, bool, error) {
	if cell.X < 0 || cell.Y < 0 || cell.X >= maze.Cols || cell.Y >= maze.Rows {
		return path, false, errors.New("invalid coords")
	}

	if cell == end {
		path[cell.Y][cell.X] = visited
		return path, true, nil
	}

	if path[cell.Y][cell.X] != unvisited {
		return path, false, nil
	}

	path[cell.Y][cell.X] = visited

	for _, neighbour := range getNeighbours(cell, maze) {
		path, res, _ := findPath(path, maze, neighbour, end)
		if res {
			path[cell.Y][cell.X] = visited
			return path, true, nil
		}
	}
	path[cell.Y][cell.X] = deadend

	return path, false, nil
}

func getNeighbours(cell domain.Coords, maze domain.Maze) []domain.Coords {
	var neighbours []domain.Coords

	if cell.Y != 0 && !maze.BottomWalls[cell.Y-1][cell.X] {
		neighbours = append(neighbours, domain.Coords{Y: cell.Y - 1, X: cell.X})
	}
	if cell.X != maze.Cols-1 && !maze.RightWalls[cell.Y][cell.X] {
		neighbours = append(neighbours, domain.Coords{Y: cell.Y, X: cell.X + 1})
	}
	if cell.Y != maze.Rows-1 && !maze.BottomWalls[cell.Y][cell.X] {
		neighbours = append(neighbours, domain.Coords{Y: cell.Y + 1, X: cell.X})
	}
	if cell.X != 0 && !maze.RightWalls[cell.Y][cell.X-1] {
		neighbours = append(neighbours, domain.Coords{Y: cell.Y, X: cell.X - 1})
	}

	return neighbours
}

func toBooleanPath(path [][]int) [][]bool {
	res := make([][]bool, len(path))
	for i := range res {
		res[i] = make([]bool, len(path[0]))
	}

	for i, row := range path {
		for j, col := range row {
			if col == visited {
				res[i][j] = true
			}
		}
	}

	return res
}
