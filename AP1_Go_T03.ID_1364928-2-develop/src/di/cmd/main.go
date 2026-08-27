package main

import (
	"net/http"
	"tictactoe/datasource"
	"tictactoe/web"

	"go.uber.org/fx"
)

func startServer(handler *web.Handler) {
	http.HandleFunc("POST /game", handler.CreateGame)
	http.HandleFunc("POST /game/{id}", handler.HandleGame)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8000", nil)
}

func main() {
	fx.New(
		fx.Provide(
			datasource.NewGameStorage,
			datasource.NewGameRepository,
			datasource.NewGameService,
			web.NewGameHandler,
		),
		fx.Invoke(startServer),
	).Run()
}
