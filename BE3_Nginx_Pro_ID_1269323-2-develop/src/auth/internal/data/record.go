package data

import "github.com/jackc/pgx/v5/pgtype"

type userRecord struct {
	Id           pgtype.UUID `db:"id"`
	Login        string      `db:"login"`
	Name         string      `db:"name"`
	Surname      string      `db:"surname"`
	Phonenumber  string      `db:"phonenumber"`
	PasswordHash []byte      `db:"passhash"`
	Email        string      `db:"email"`
	Role         string      `db:"role"`
}
