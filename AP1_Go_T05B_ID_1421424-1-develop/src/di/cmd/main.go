package main

import (
	"net/http"
	"tictactoe/datasource"
	"tictactoe/web"

	"go.uber.org/fx"
)

func startServer(gameHandler *web.Handler, authHandler *web.AuthHandler, userHandler *web.UserHandler, authenticator *web.UserAuthenticator) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/signup", authHandler.SignUp)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/auth/access", authHandler.RefreshAccessToken)
	mux.HandleFunc("/auth/refresh", authHandler.RefreshRefreshToken)

	mux.Handle("POST /game", authenticator.Authenticate(http.HandlerFunc(gameHandler.CreateGame)))
	mux.Handle("GET /games/available", authenticator.Authenticate(http.HandlerFunc(gameHandler.GetAvailableGames)))
	mux.Handle("GET /games/completed", authenticator.Authenticate(http.HandlerFunc(gameHandler.GetCompletedGames)))
	mux.Handle("GET /games/leaderboard", authenticator.Authenticate(http.HandlerFunc(gameHandler.GetTopPlayers)))
	mux.Handle("POST /game/{id}", authenticator.Authenticate(http.HandlerFunc(gameHandler.HandleGame)))
	mux.Handle("GET /game/{id}", authenticator.Authenticate(http.HandlerFunc(gameHandler.GetGame)))
	mux.Handle("POST /game/{id}/join", authenticator.Authenticate(http.HandlerFunc(gameHandler.JoinGame)))
	mux.Handle("GET /users/{id}", authenticator.Authenticate(http.HandlerFunc(userHandler.GetUser)))
	mux.Handle("GET /users/me", authenticator.Authenticate(http.HandlerFunc(userHandler.GetCurrentUser)))

	return http.ListenAndServe(":8000", mux)
}

func main() {
	fx.New(
		fx.Provide(
			datasource.NewGameRepository,
			datasource.NewGameService,
			datasource.NewPostgresPool,
			datasource.NewUserRepository,
			datasource.NewUserService,
			datasource.NewAuthService,
			datasource.NewJwtProvider,
			datasource.NewRefreshTokenRepository,
			web.NewGameHandler,
			web.NewAuthHandler,
			web.NewUserHandler,
			web.NewUserAuthenticator,
		),
		fx.Invoke(startServer),
	).Run()
}
