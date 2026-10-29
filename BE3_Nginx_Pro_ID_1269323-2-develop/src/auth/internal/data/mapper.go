package data

import (
	"auth/internal/auth"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

func toPgUUID(id uuid.UUID) pgtype.UUID {
	res := pgtype.UUID{Bytes: id, Valid: true}
	if id == uuid.Nil {
		res.Valid = false
	}

	return res
}

func toUser(record userRecord) *auth.User {
	id, err := uuid.FromBytes(record.Id.Bytes[:])
	if err != nil {
		id = uuid.Nil
	}

	return &auth.User{
		Id:           id,
		Login:        record.Login,
		Name:         record.Name,
		Surname:      record.Surname,
		Phonenumber:  record.Phonenumber,
		PasswordHash: record.PasswordHash,
		Email:        record.Email,
		Role:         record.Role,
	}
}
