package mazefile

import (
	"bytes"
	"strings"
	"testing"

	"maze/domain"
	"maze/domain/eller"
)

// Проверка сквозного сохранения и чтения (Round-Trip)
func TestFormatterRoundTrip(t *testing.T) {
	original, err := eller.Generate(10, 15)
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	var buf bytes.Buffer
	formatter := NewFormatter()
	if err := formatter.Format(&buf, original); err != nil {
		t.Fatalf("Format() unexpected error: %v", err)
	}

	parser := NewParser()
	parsed, err := parser.Parse(&buf)
	if err != nil {
		t.Fatalf("Parse() formatted output failed: %v", err)
	}

	if parsed.Rows != original.Rows || parsed.Cols != original.Cols {
		t.Fatalf("Parsed dimensions %dx%d != original %dx%d", parsed.Rows, parsed.Cols, original.Rows, original.Cols)
	}

	for r := 0; r < original.Rows; r++ {
		for c := 0; c < original.Cols; c++ {
			if parsed.RightWalls[r][c] != original.RightWalls[r][c] {
				t.Fatalf("RightWalls mismatch at [%d,%d]", r, c)
			}
			if parsed.BottomWalls[r][c] != original.BottomWalls[r][c] {
				t.Fatalf("BottomWalls mismatch at [%d,%d]", r, c)
			}
		}
	}
}

// Проверка валидации и возврата ошибок при некорректных входных данных
func TestFormatterRejectsInvalid(t *testing.T) {
	formatter := NewFormatter()

	// Передача nil вместо io.Writer
	t.Run("nil writer", func(t *testing.T) {
		maze := domain.NewMaze(2, 2, [][]bool{{true, true}, {true, true}}, [][]bool{{true, true}, {true, true}})
		if err := formatter.Format(nil, maze); err == nil {
			t.Fatal("Format(nil, ...) expected error, got nil")
		}
	})

	// Размеры выходят за границы 1..50
	t.Run("invalid dimensions", func(t *testing.T) {
		var buf bytes.Buffer
		invalidMaze := domain.Maze{Rows: 0, Cols: 10}
		if err := formatter.Format(&buf, invalidMaze); err == nil {
			t.Fatal("Format() with 0 rows expected error, got nil")
		}

		tooBigMaze := domain.Maze{Rows: 51, Cols: 10}
		if err := formatter.Format(&buf, tooBigMaze); err == nil {
			t.Fatal("Format() with 51 rows expected error, got nil")
		}
	})

	// Несовпадение матриц стен с размерами лабиринта
	t.Run("matrix size mismatch", func(t *testing.T) {
		var buf bytes.Buffer
		mismatched := domain.Maze{
			Rows:        3,
			Cols:        3,
			RightWalls:  [][]bool{{true, true, true}},
			BottomWalls: [][]bool{{true, true, true}},
		}
		if err := formatter.Format(&buf, mismatched); err == nil {
			t.Fatal("Format() with mismatched matrix expected error, got nil")
		}
	})
}

// Проверка точного формата вывода
func TestFormatterOutputFormat(t *testing.T) {
	maze := domain.NewMaze(
		2, 3,
		[][]bool{
			{false, true, true},
			{true, false, true},
		},
		[][]bool{
			{false, true, false},
			{true, true, true},
		},
	)

	var buf bytes.Buffer
	if err := NewFormatter().Format(&buf, maze); err != nil {
		t.Fatalf("Format() unexpected error: %v", err)
	}

	expected := "2 3\n0 1 1\n1 0 1\n\n0 1 0\n1 1 1\n"
	if strings.ReplaceAll(buf.String(), "\r\n", "\n") != expected {
		t.Fatalf("Formatted output does not match expected.\nGot:\n%s\nWant:\n%s", buf.String(), expected)
	}
}
