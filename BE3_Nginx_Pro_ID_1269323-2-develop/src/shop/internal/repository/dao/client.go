package dao

import "time"

type Client struct {
	ID               string
	ClientName       string
	ClientSurname    string
	Birthday         time.Time
	Gender           string
	RegistrationDate time.Time
	AddressID        string
}
