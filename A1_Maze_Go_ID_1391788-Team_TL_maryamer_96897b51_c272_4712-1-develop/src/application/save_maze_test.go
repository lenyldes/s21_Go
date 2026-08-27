package application

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"maze/domain"
)

type formatterStub struct {
	err error
}

func (stub formatterStub) Format(writer io.Writer, maze domain.Maze) error {
	if stub.err != nil {
		return stub.err
	}
	_, err := writer.Write([]byte("saved"))
	return err
}

func TestSaveMazeSavesCurrentMaze(t *testing.T) {
	session := NewMazeSession()
	session.setCurrent(domain.NewMaze(1, 1, [][]bool{{true}}, [][]bool{{true}}))

	var buf bytes.Buffer
	saveMaze := NewSaveMaze(formatterStub{}, session)
	if err := saveMaze.Execute(&buf); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}
	if buf.String() != "saved" {
		t.Fatalf("Execute() written = %q, want %q", buf.String(), "saved")
	}
}

func TestSaveMazeFailsWhenNoMazeInSession(t *testing.T) {
	session := NewMazeSession()
	var buf bytes.Buffer
	saveMaze := NewSaveMaze(formatterStub{}, session)
	if err := saveMaze.Execute(&buf); err == nil {
		t.Fatal("Execute() error = nil, want error when session has no maze")
	}
}

func TestSaveMazeReturnsFormatterError(t *testing.T) {
	session := NewMazeSession()
	session.setCurrent(domain.NewMaze(1, 1, [][]bool{{true}}, [][]bool{{true}}))

	formatErr := errors.New("format error")
	saveMaze := NewSaveMaze(formatterStub{err: formatErr}, session)
	var buf bytes.Buffer
	if err := saveMaze.Execute(&buf); !errors.Is(err, formatErr) {
		t.Fatalf("Execute() error = %v, want %v", err, formatErr)
	}
}
