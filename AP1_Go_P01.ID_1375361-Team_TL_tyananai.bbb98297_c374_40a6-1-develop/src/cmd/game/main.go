package main

import (
	"rogue/internal/representation"

	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	view := representation.NewView()
	view.Render()
	for !view.IsGameOver() {
		input := view.HandleInput()
		view.Update(input)
		view.Render()
	}
}
