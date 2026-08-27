package application

import (
	"io"

	"maze/domain"
)

type MazeParser interface {
	Parse(io.Reader) (domain.Maze, error)
}

type LoadMaze struct {
	parser  MazeParser
	session *MazeSession
}

func NewLoadMaze(parser MazeParser, session *MazeSession) LoadMaze {
	return LoadMaze{
		parser:  parser,
		session: session,
	}
}

func (useCase LoadMaze) Execute(reader io.Reader) error {
	maze, err := useCase.parser.Parse(reader)
	if err != nil {
		return err
	}

	useCase.session.setCurrent(maze)
	return nil
}
