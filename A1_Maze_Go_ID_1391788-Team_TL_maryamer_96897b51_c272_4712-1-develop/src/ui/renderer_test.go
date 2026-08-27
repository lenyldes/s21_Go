package ui

import (
	"image/color"
	"testing"

	"maze/domain"
)

func TestRenderMazeSizeAndWalls(t *testing.T) {
	maze := domain.NewMaze(1, 1, [][]bool{{true}}, [][]bool{{true}})
	result := RenderMaze(maze)
	if result.Bounds().Dx() != 500 || result.Bounds().Dy() != 500 {
		t.Fatalf("RenderMaze() size = %v, want 500x500", result.Bounds())
	}
	black := color.RGBA{A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if got := result.At(0, 250); got != black {
		t.Fatalf("left wall color = %v, want %v", got, black)
	}
	if got := result.At(250, 250); got != white {
		t.Fatalf("cell color = %v, want %v", got, white)
	}
	if got := result.At(499, 250); got != black {
		t.Fatalf("right wall color = %v, want %v", got, black)
	}
}
