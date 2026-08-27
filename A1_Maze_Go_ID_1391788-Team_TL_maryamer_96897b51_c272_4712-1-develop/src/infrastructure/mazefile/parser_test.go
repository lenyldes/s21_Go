package mazefile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParserValidMaze(t *testing.T) {
	input := `2 3
0 1 1
1 0 1

0 1 0
1 1 1`
	maze, err := NewParser().Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if maze.Rows != 2 || maze.Cols != 3 {
		t.Fatalf("Parse() dimensions = %dx%d, want 2x3", maze.Rows, maze.Cols)
	}
	if !maze.RightWalls[0][1] || maze.BottomWalls[0][0] {
		t.Fatal("Parse() returned incorrect wall matrices")
	}
}

func TestParserRejectsInvalidFiles(t *testing.T) {
	tests := map[string]string{
		"missing dimensions":  "",
		"dimension too large": "51 1",
		"invalid wall value":  "1 1\n2\n1",
		"missing wall":        "1 1\n1",
		"extra data":          "1 1\n1\n1\n0",
		"open right edge":     "1 1\n0\n1",
		"open bottom edge":    "1 1\n1\n0",
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := NewParser().Parse(strings.NewReader(input)); err == nil {
				t.Fatal("Parse() error = nil, want validation error")
			}
		})
	}
}

func TestExampleFilesAreValid(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "examples", "maze_*.txt"))
	if err != nil {
		t.Fatalf("Glob() unexpected error: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no example maze files found")
	}
	for _, path := range files {
		t.Run(filepath.Base(path), func(t *testing.T) {
			file, openErr := os.Open(path)
			if openErr != nil {
				t.Fatalf("Open() unexpected error: %v", openErr)
			}
			defer file.Close()
			if _, parseErr := NewParser().Parse(file); parseErr != nil {
				t.Fatalf("Parse() unexpected error: %v", parseErr)
			}
		})
	}
}
