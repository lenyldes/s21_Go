package ui

import (
	"fmt"
	"image/color"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"maze/application"
	"maze/domain"
)

const initialStatusMessage = "Выберите файл лабиринта или сгенерируйте новый"

type Window struct {
	window          fyne.Window
	loadMaze        application.LoadMaze
	saveMaze        application.SaveMaze
	generateMaze    application.GenerateMaze
	coordsProcessor application.CoordsProcessor
	session         *application.MazeSession
	preview         *fyne.Container
	frame           *canvas.Rectangle
	status          *widget.Label
}

func NewWindow(
	app fyne.App,
	loadMaze application.LoadMaze,
	saveMaze application.SaveMaze,
	generateMaze application.GenerateMaze,
	coordsProcessor application.CoordsProcessor,
	session *application.MazeSession,
) *Window {
	mainWindow := app.NewWindow("Maze")
	emptyCanvas := canvas.NewRectangle(color.White)
	emptyCanvas.SetMinSize(fyne.NewSize(mazeCanvasSize, mazeCanvasSize))
	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeColor = color.Gray{Y: 96}
	frame.StrokeWidth = 2
	frame.SetMinSize(fyne.NewSize(mazeCanvasSize, mazeCanvasSize))
	view := &Window{
		window:          mainWindow,
		loadMaze:        loadMaze,
		saveMaze:        saveMaze,
		generateMaze:    generateMaze,
		coordsProcessor: coordsProcessor,
		session:         session,
		preview:         container.NewStack(emptyCanvas, frame),
		frame:           frame,
		status:          widget.NewLabel(initialStatusMessage),
	}
	view.preview.Resize(fyne.NewSize(mazeCanvasSize, mazeCanvasSize))
	workspace := container.NewCenter(view.preview)

	loadButton := widget.NewButton("Загрузить", view.openFile)
	saveButton := widget.NewButton("Сохранить", view.saveFile)
	generateButton := widget.NewButton("Сгенерировать", view.promptGenerate)
	findPathButton := widget.NewButton("Найти путь", view.solve)
	clearButton := widget.NewButton("Очистить", view.clear)
	buttons := container.NewGridWithColumns(
		5,
		loadButton,
		saveButton,
		generateButton,
		findPathButton,
		clearButton,
	)
	controls := container.NewVBox(
		widget.NewLabelWithStyle("Лабиринт", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		buttons,
		view.status,
	)
	content := container.NewVBox(controls, workspace)
	mainWindow.SetContent(container.NewPadded(content))
	mainWindow.Resize(fyne.NewSize(750, 750))
	mainWindow.SetFixedSize(true)
	return view
}

func (view *Window) ShowAndRun() {
	view.window.ShowAndRun()
}

func (view *Window) openFile() {
	fileDialog := dialog.NewFileOpen(view.fileSelected, view.window)
	fileDialog.Resize(fyne.NewSize(650, 550))
	fileDialog.Show()
}

func (view *Window) fileSelected(reader fyne.URIReadCloser, dialogErr error) {
	if dialogErr != nil {
		dialog.ShowError(dialogErr, view.window)
		return
	}
	if reader == nil {
		return
	}
	defer closeReader(reader)

	if err := view.loadMaze.Execute(reader); err != nil {
		dialog.ShowError(fmt.Errorf("некорректный файл лабиринта: %w", err), view.window)
		return
	}
	maze, exists := view.session.Current()
	if !exists {
		dialog.ShowError(fmt.Errorf("лабиринт не загружен"), view.window)
		return
	}

	view.showMaze(maze, filepath.Base(reader.URI().Path()))
}

func (view *Window) promptGenerate() {
	rowsEntry := widget.NewEntry()
	rowsEntry.SetPlaceHolder("1..50")
	colsEntry := widget.NewEntry()
	colsEntry.SetPlaceHolder("1..50")

	formItems := []*widget.FormItem{
		widget.NewFormItem("Строки:", rowsEntry),
		widget.NewFormItem("Столбцы:", colsEntry),
	}

	formDialog := dialog.NewForm(
		"Параметры генерации",
		"Сгенерировать",
		"Отмена",
		formItems,
		func(confirm bool) {
			if !confirm {
				return
			}
			rows, errR := strconv.Atoi(strings.TrimSpace(rowsEntry.Text))
			cols, errC := strconv.Atoi(strings.TrimSpace(colsEntry.Text))
			if errR != nil || errC != nil || rows < 1 || rows > 50 || cols < 1 || cols > 50 {
				dialog.ShowError(fmt.Errorf("размеры лабиринта должны быть числами от 1 до 50"), view.window)
				return
			}

			if err := view.generateMaze.Execute(rows, cols); err != nil {
				dialog.ShowError(fmt.Errorf("ошибка генерации: %w", err), view.window)
				return
			}

			maze, exists := view.session.Current()
			if !exists {
				dialog.ShowError(fmt.Errorf("лабиринт не сгенерирован"), view.window)
				return
			}
			view.showMaze(maze, "Сгенерированный лабиринт")
		},
		view.window,
	)
	formDialog.Resize(fyne.NewSize(350, 180))
	formDialog.Show()
}

func (view *Window) saveFile() {
	if _, exists := view.session.Current(); !exists {
		dialog.ShowInformation("Сохранение", "Нет активного лабиринта для сохранения", view.window)
		return
	}

	saveDialog := dialog.NewFileSave(view.fileSaved, view.window)
	saveDialog.Resize(fyne.NewSize(650, 550))
	saveDialog.SetFileName("maze.txt")
	saveDialog.Show()
}

func (view *Window) fileSaved(writer fyne.URIWriteCloser, dialogErr error) {
	if dialogErr != nil {
		dialog.ShowError(dialogErr, view.window)
		return
	}
	if writer == nil {
		return
	}
	defer closeWriter(writer)

	if err := view.saveMaze.Execute(writer); err != nil {
		dialog.ShowError(fmt.Errorf("ошибка сохранения: %w", err), view.window)
		return
	}
	dialog.ShowInformation("Сохранение", "Лабиринт успешно сохранен", view.window)
}

func (view *Window) showMaze(maze domain.Maze, title string) {
	mazeImage := canvas.NewImageFromImage(RenderMaze(maze))
	mazeImage.FillMode = canvas.ImageFillOriginal
	mazeImage.SetMinSize(fyne.NewSize(mazeCanvasSize, mazeCanvasSize))
	mazeImage.Resize(fyne.NewSize(mazeCanvasSize, mazeCanvasSize))
	view.preview.Objects = []fyne.CanvasObject{mazeImage, view.frame}
	view.preview.Refresh()
	view.status.SetText(fmt.Sprintf("%s — %d×%d", title, maze.Rows, maze.Cols))
}

func closeReader(reader io.Closer) {
	_ = reader.Close()
}

func closeWriter(writer io.Closer) {
	_ = writer.Close()
}

func (view *Window) solve() {
	maze, exists := view.session.Current()
	if !exists {
		dialog.ShowError(fmt.Errorf("лабиринт не загружен"), view.window)
		return
	}

	exit1StringCoords := widget.NewEntry()
	exit2StringCoords := widget.NewEntry()
	formWidgets := []*widget.FormItem{
		&widget.FormItem{Text: "Выход 1:", Widget: exit1StringCoords},
		&widget.FormItem{Text: "Выход 2:", Widget: exit2StringCoords},
	}

	dialog.ShowForm(
		"Введите координаты выходов в формате X:Y",
		"Ок",
		"Отмена",
		formWidgets,
		func(confirm bool) {
			if !confirm {
				return
			}

			exit1, exit2, err := view.coordsProcessor.GetCoords(exit1StringCoords.Text, exit2StringCoords.Text, maze.Rows, maze.Cols)
			if err != nil {
				dialog.ShowError(err, view.window)
				return
			}

			if err := view.session.Solve(exit1, exit2); err != nil {
				dialog.ShowError(err, view.window)
				return
			}

			path := RenderPath(maze.Rows, maze.Cols, toLineSlice(view.session.FindPath(), maze, view.session.Start()))
			view.preview.Objects = []fyne.CanvasObject{view.preview.Objects[0], view.preview.Objects[1], path}
			view.preview.Refresh()
		},
		view.window,
	)
}

func (view *Window) clear() {
	view.session.Clear()
	emptyCanvas := canvas.NewRectangle(color.White)
	emptyCanvas.SetMinSize(fyne.NewSize(mazeCanvasSize, mazeCanvasSize))
	view.preview.Objects = []fyne.CanvasObject{emptyCanvas, view.frame}
	view.preview.Refresh()
	view.status.SetText(initialStatusMessage)
}
