package main

import (
	"fyne.io/fyne/v2/app"

	"maze/application"
	"maze/domain/eller"
	"maze/infrastructure/coords"
	"maze/infrastructure/mazefile"
	"maze/ui"
)

func main() {
	mazeApp := app.NewWithID("school21.maze")
	session := application.NewMazeSession()
	loader := application.NewLoadMaze(mazefile.NewParser(), session)
	saver := application.NewSaveMaze(mazefile.NewFormatter(), session)
	generator := application.NewGenerateMaze(eller.NewGenerator(), session)
	coordsProcessor := application.NewCoordsProcessor(coords.NewParser())
	ui.NewWindow(mazeApp, loader, saver, generator, coordsProcessor, session).ShowAndRun()
}
