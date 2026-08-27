package application

import (
	"errors"
	"io"
	"strings"
	"testing"

	"maze/domain"
)

type parserStub struct {
	maze domain.Maze
	err  error
}

func (stub parserStub) Parse(io.Reader) (domain.Maze, error) {
	return stub.maze, stub.err
}

func TestLoadMazeStoresCurrentMaze(t *testing.T) {
	session := NewMazeSession()
	if _, exists := session.Current(); exists {
		t.Fatal("Current() exists = true before loading")
	}

	expected := domain.NewMaze(1, 1, [][]bool{{true}}, [][]bool{{true}})
	loadMaze := NewLoadMaze(parserStub{maze: expected}, session)
	if err := loadMaze.Execute(strings.NewReader("maze")); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	current, exists := session.Current()
	if !exists {
		t.Fatal("Current() exists = false after loading")
	}
	if current.Rows != expected.Rows || current.Cols != expected.Cols {
		t.Fatalf("Current() dimensions = %dx%d, want %dx%d", current.Rows, current.Cols, expected.Rows, expected.Cols)
	}
}

func TestLoadMazeKeepsCurrentMazeAfterParserError(t *testing.T) {
	session := NewMazeSession()
	expected := domain.NewMaze(1, 1, [][]bool{{true}}, [][]bool{{true}})
	if err := NewLoadMaze(parserStub{maze: expected}, session).Execute(strings.NewReader("maze")); err != nil {
		t.Fatalf("first Execute() unexpected error: %v", err)
	}

	parseErr := errors.New("parse error")
	if err := NewLoadMaze(parserStub{err: parseErr}, session).Execute(strings.NewReader("broken")); !errors.Is(err, parseErr) {
		t.Fatalf("second Execute() error = %v, want %v", err, parseErr)
	}

	current, exists := session.Current()
	if !exists || current.Rows != expected.Rows || current.Cols != expected.Cols {
		t.Fatal("Current() changed after parser error")
	}
}

func TestMazeSessionClear(t *testing.T) {
	session := NewMazeSession()
	maze := domain.NewMaze(2, 2, [][]bool{{false, true}, {false, true}}, [][]bool{{true, true}, {true, true}})
	session.setCurrent(maze)

	start := domain.Coords{X: 0, Y: 0}
	end := domain.Coords{X: 1, Y: 0}
	if err := session.Solve(start, end); err != nil {
		t.Fatalf("Solve() unexpected error: %v", err)
	}

	if _, exists := session.Current(); !exists {
		t.Fatal("Current() exists = false before Clear()")
	}
	if len(session.FindPath()) == 0 {
		t.Fatal("FindPath() is empty before Clear()")
	}

	session.Clear()

	if _, exists := session.Current(); exists {
		t.Fatal("Current() exists = true after Clear()")
	}
	if len(session.FindPath()) != 0 {
		t.Fatal("FindPath() is not empty after Clear()")
	}
	if session.Start() != (domain.Coords{}) {
		t.Fatalf("Start() = %v after Clear(), want zero coords", session.Start())
	}
}
