package ui

import (
	"maze/domain"
)

func toLineSlice(path [][]bool, maze domain.Maze, start domain.Coords) []Line {
	var res []Line

	cur := start

	change := true

	for change {
		change = false
		startLine := cur

		path[cur.Y][cur.X] = false

		// right
		for cur.X+1 < maze.Cols && path[cur.Y][cur.X+1] && !maze.RightWalls[cur.Y][cur.X] {
			cur.X += 1
			path[cur.Y][cur.X] = false
		}

		if cur != startLine {
			res = append(res, Line{Start: startLine, End: cur})
			change = true
			continue
		}

		//left
		for cur.X-1 >= 0 && path[cur.Y][cur.X-1] && !maze.RightWalls[cur.Y][cur.X-1] {
			cur.X -= 1
			path[cur.Y][cur.X] = false
		}

		if cur != startLine {
			res = append(res, Line{Start: startLine, End: cur})
			change = true
			continue
		}

		// down
		for cur.Y+1 < maze.Rows && path[cur.Y+1][cur.X] && !maze.BottomWalls[cur.Y][cur.X] {
			cur.Y += 1
			path[cur.Y][cur.X] = false
		}
		if cur != startLine {
			res = append(res, Line{Start: startLine, End: cur})
			change = true
			continue
		}

		// up
		for cur.Y-1 >= 0 && path[cur.Y-1][cur.X] && !maze.BottomWalls[cur.Y-1][cur.X] {
			cur.Y -= 1
			path[cur.Y][cur.X] = false
		}

		if cur != startLine {
			res = append(res, Line{Start: startLine, End: cur})
			change = true
		}

	}

	return res
}
