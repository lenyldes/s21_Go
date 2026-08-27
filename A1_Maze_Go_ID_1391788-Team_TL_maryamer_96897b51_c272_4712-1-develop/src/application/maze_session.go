package application

import (
	"maze/domain"
)

type MazeSession struct {
	current    domain.Maze
	hasCurrent bool
	findPath   [][]bool
	start      domain.Coords
	end        domain.Coords
}

func NewMazeSession() *MazeSession {
	return &MazeSession{}
}

func (session *MazeSession) Current() (domain.Maze, bool) {
	if session == nil || !session.hasCurrent {
		return domain.Maze{}, false
	}

	return session.current, true
}

func (session *MazeSession) setCurrent(maze domain.Maze) {
	session.current = maze
	session.hasCurrent = true
}

func (session *MazeSession) Solve(start domain.Coords, end domain.Coords) error {
	session.start = start
	session.end = end

	path, err := Solve(session.current, start, end)
	if err != nil {
		return err
	}

	session.findPath = path
	return nil
}

func (session *MazeSession) FindPath() [][]bool {
	var res [][]bool
	for i, row := range session.findPath {
		res = append(res, make([]bool, session.current.Cols))
		copy(res[i], row)
	}

	return res
}

func (session *MazeSession) Start() domain.Coords {
	return session.start
}

func (session *MazeSession) Clear() {
	if session == nil {
		return
	}
	session.current = domain.Maze{}
	session.hasCurrent = false
	session.findPath = nil
	session.start = domain.Coords{}
	session.end = domain.Coords{}
}
