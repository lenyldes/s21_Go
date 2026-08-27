package application

import (
	"errors"
	"io"

	"maze/domain"
)

type MazeFormatter interface {
	Format(io.Writer, domain.Maze) error
}

type SaveMaze struct {
	formatter MazeFormatter
	session   *MazeSession
}

func NewSaveMaze(formatter MazeFormatter, session *MazeSession) SaveMaze {
	return SaveMaze{
		formatter: formatter,
		session:   session,
	}
}

func (useCase SaveMaze) Execute(writer io.Writer) error {
	maze, exists := useCase.session.Current()
	if !exists {
		return errors.New("лабиринт не задан")
	}

	return useCase.formatter.Format(writer, maze)
}
