package internal

import "time"

type BoardWeb [3][3]string

type GameRequest struct {
	Board BoardWeb `json:"board"`
}

type GameResponse struct {
	UUID       string    `json:"uuid"`
	Board      BoardWeb  `json:"board"`
	GameOver   bool      `json:"game_over"`
	Phase      string    `json:"phase"`
	PlayerUUID string    `json:"player_uuid"`
	PlayerX    string    `json:"player_x"`
	PlayerO    string    `json:"player_o"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserResponse struct {
	UUID string `json:"uuid"`
}

type UserInfoResponse struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`
}

type TopPlayerResponse struct {
	UUID     string  `json:"uuid"`
	Login    string  `json:"login"`
	WinRatio float64 `json:"win_ratio"`
}

type contextKey string

const userUUIDKey contextKey = "uuid"

type CreateGameRequest struct {
	Mode string `json:"game_mode"`
}
