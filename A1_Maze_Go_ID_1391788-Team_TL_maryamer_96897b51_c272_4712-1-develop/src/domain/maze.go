package domain

type Maze struct {
	Rows        int
	Cols        int
	RightWalls  [][]bool
	BottomWalls [][]bool
	Start       Coords
	End         Coords
}

func NewMaze(rows, cols int, rightWalls, bottomWalls [][]bool) Maze {
	return Maze{
		Rows:        rows,
		Cols:        cols,
		RightWalls:  rightWalls,
		BottomWalls: bottomWalls,
	}
}
