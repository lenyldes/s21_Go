package internal

import "time"

type BoardEntity [9]string

type GameEntity struct {
	UUID       string      `db:"uuid"`
	Board      BoardEntity `db:"board"`
	Phase      string      `db:"phase"`
	PlayerUUID string      `db:"player_uuid"`
	PlayerX    string      `db:"player_x"`
	PlayerO    string      `db:"player_o"`
	CreatedAt  time.Time   `db:"created_at"`
}

type UserEntity struct {
	UUID     string `db:"uuid"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
