package ui

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"maze/domain"
)

const (
	mazeCanvasSize = 500
	wallThickness  = 2
)

func RenderMaze(maze domain.Maze) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, mazeCanvasSize, mazeCanvasSize))
	fill(result, color.White)

	drawVertical(result, 0)
	drawHorizontal(result, 0)
	for row := 0; row < maze.Rows; row++ {
		for col := 0; col < maze.Cols; col++ {
			if maze.RightWalls[row][col] {
				x := boundary(col+1, maze.Cols)
				y1 := boundary(row, maze.Rows)
				y2 := boundary(row+1, maze.Rows)
				drawVerticalRange(result, x, y1, y2)
			}
			if maze.BottomWalls[row][col] {
				y := boundary(row+1, maze.Rows)
				x1 := boundary(col, maze.Cols)
				x2 := boundary(col+1, maze.Cols)
				drawHorizontalRange(result, y, x1, x2)
			}
		}
	}
	return result
}

func boundary(index, count int) int {
	return int(math.Round(float64(index) * mazeCanvasSize / float64(count)))
}

func fill(target *image.RGBA, value color.Color) {
	for y := 0; y < mazeCanvasSize; y++ {
		for x := 0; x < mazeCanvasSize; x++ {
			target.Set(x, y, value)
		}
	}
}

func drawVertical(target *image.RGBA, x int) {
	drawVerticalRange(target, x, 0, mazeCanvasSize)
}

func drawHorizontal(target *image.RGBA, y int) {
	drawHorizontalRange(target, y, 0, mazeCanvasSize)
}

func drawVerticalRange(target *image.RGBA, x, start, end int) {
	startX := x
	if x == mazeCanvasSize {
		startX = mazeCanvasSize - wallThickness
	}
	for y := start; y < end; y++ {
		for offset := 0; offset < wallThickness; offset++ {
			target.Set(startX+offset, y, color.Black)
		}
	}
}

func drawHorizontalRange(target *image.RGBA, y, start, end int) {
	startY := y
	if y == mazeCanvasSize {
		startY = mazeCanvasSize - wallThickness
	}
	for x := start; x < end; x++ {
		for offset := 0; offset < wallThickness; offset++ {
			target.Set(x, startY+offset, color.Black)
		}
	}
}

func RenderPath(rows, cols int, path []Line) *fyne.Container {
	cellSize := float32(mazeCanvasSize) / float32(rows)

	cn := container.NewWithoutLayout()
	for _, linePath := range path {
		line := canvas.NewLine(color.NRGBA{67, 245, 39, 255})
		line.StrokeWidth = 2
		line.Position1 = fyne.Position{X: float32(linePath.Start.X)*cellSize + cellSize/2 + wallThickness/2, Y: float32(linePath.Start.Y)*cellSize + cellSize/2 + wallThickness/2}
		line.Position2 = fyne.Position{X: float32(linePath.End.X)*cellSize + cellSize/2 + wallThickness/2, Y: float32(linePath.End.Y)*cellSize + cellSize/2 + wallThickness/2}
		cn.Add(line)
	}

	cn.Refresh()
	return cn
}
