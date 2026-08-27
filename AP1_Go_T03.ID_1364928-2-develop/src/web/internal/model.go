package internal

type BoardWeb [3][3]int

type GameRequest struct {
	Board BoardWeb `json:"board"`
}

type GameResponse struct {
	UUID     string   `json:"uuid"`
	Board    BoardWeb `json:"board"`
	GameOver bool     `json:"game_over"`
}
