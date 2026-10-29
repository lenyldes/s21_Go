package web

import (
	"net/http"
	"tictactoe/domain"
	"tictactoe/web/internal"
)

type Handler struct {
	handler *internal.GameHandler
}

func NewGameHandler(service domain.GameService, repo domain.GameRepository, userService domain.UserService) *Handler {
	return &Handler{
		handler: internal.NewGameHandler(service, repo, userService),
	}
}

func (h *Handler) HandleGame(w http.ResponseWriter, r *http.Request) {
	h.handler.HandleGame(w, r)
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	h.handler.CreateGame(w, r)
}

func (h *Handler) JoinGame(w http.ResponseWriter, r *http.Request) {
	h.handler.JoinGame(w, r)
}

func (h *Handler) GetAvailableGames(w http.ResponseWriter, r *http.Request) {
	h.handler.GetAvailableGames(w, r)
}

func (h *Handler) GetCompletedGames(w http.ResponseWriter, r *http.Request) {
	h.handler.GetCompletedGames(w, r)
}

func (h *Handler) GetTopPlayers(w http.ResponseWriter, r *http.Request) {
	h.handler.GetTopPlayers(w, r)
}

func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	h.handler.GetGame(w, r)
}

type AuthHandler struct {
	handler *internal.AuthHandler
}

func NewAuthHandler(service domain.AuthService) *AuthHandler {
	return &AuthHandler{
		handler: internal.NewAuthHandler(service),
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	h.handler.SignUp(w, r)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.handler.Login(w, r)
}

func (h *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	h.handler.RefreshAccessToken(w, r)
}

func (h *AuthHandler) RefreshRefreshToken(w http.ResponseWriter, r *http.Request) {
	h.handler.RefreshRefreshToken(w, r)
}

type UserHandler struct {
	handler *internal.UserHandler
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{
		handler: internal.NewUserHandler(service),
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.handler.GetUser(w, r)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	h.handler.GetCurrentUser(w, r)
}

type UserAuthenticator struct {
	handler *internal.UserAuthenticator
}

func NewUserAuthenticator(provider domain.JwtProvider) *UserAuthenticator {
	return &UserAuthenticator{
		handler: internal.NewUserAuthenticator(provider),
	}
}

func (h *UserAuthenticator) Authenticate(next http.Handler) http.Handler {
	return h.handler.Authenticate(next)
}
