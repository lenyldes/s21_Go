package application

import (
	"maze/domain"
)

type MazeGenerator interface {
	Generate(rows, cols int) (domain.Maze, error)
}

type GenerateMaze struct {
	generator MazeGenerator
	session   *MazeSession
}

func NewGenerateMaze(generator MazeGenerator, session *MazeSession) GenerateMaze {
	return GenerateMaze{
		generator: generator,
		session:   session,
	}
}

func (useCase GenerateMaze) Execute(rows, cols int) error {
	maze, err := useCase.generator.Generate(rows, cols)
	if err != nil {
		return err
	}

	useCase.session.setCurrent(maze)
	return nil
}
