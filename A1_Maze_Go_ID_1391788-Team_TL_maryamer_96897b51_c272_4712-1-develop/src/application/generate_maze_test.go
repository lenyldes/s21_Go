package application

import (
	"errors"
	"testing"

	"maze/domain"
)

type generatorStub struct {
	maze domain.Maze
	err  error
}

func (stub generatorStub) Generate(rows, cols int) (domain.Maze, error) {
	return stub.maze, stub.err
}

func TestGenerateMazeStoresCurrentMaze(t *testing.T) {
	session := NewMazeSession()
	expected := domain.NewMaze(2, 2, [][]bool{{true, true}, {true, true}}, [][]bool{{true, true}, {true, true}})
	generateMaze := NewGenerateMaze(generatorStub{maze: expected}, session)

	if err := generateMaze.Execute(2, 2); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	current, exists := session.Current()
	if !exists {
		t.Fatal("Current() exists = false after generation")
	}
	if current.Rows != expected.Rows || current.Cols != expected.Cols {
		t.Fatalf("Current() dimensions = %dx%d, want %dx%d", current.Rows, current.Cols, expected.Rows, expected.Cols)
	}
}

func TestGenerateMazeKeepsCurrentMazeAfterError(t *testing.T) {
	session := NewMazeSession()
	expected := domain.NewMaze(1, 1, [][]bool{{true}}, [][]bool{{true}})
	if err := NewGenerateMaze(generatorStub{maze: expected}, session).Execute(1, 1); err != nil {
		t.Fatalf("first Execute() unexpected error: %v", err)
	}

	genErr := errors.New("generation error")
	if err := NewGenerateMaze(generatorStub{err: genErr}, session).Execute(-1, -1); !errors.Is(err, genErr) {
		t.Fatalf("second Execute() error = %v, want %v", err, genErr)
	}

	current, exists := session.Current()
	if !exists || current.Rows != expected.Rows || current.Cols != expected.Cols {
		t.Fatal("Current() changed after generator error")
	}
}
