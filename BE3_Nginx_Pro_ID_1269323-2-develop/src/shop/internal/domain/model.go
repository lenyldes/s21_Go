package domain

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

type Product struct {
	ID             string
	Name           string
	Category       string
	Price          float64
	AvailableStock int
	LastUpdateDate time.Time
	SupplierID     string
	ImageID        *string // nil, если изображение ещё не добавлено.
}

type Supplier struct {
	ID          string
	Name        string
	AddressID   string
	PhoneNumber string
}

type Image struct {
	ID    string
	Image []byte
}

type Address struct {
	ID      string
	Country string
	City    string
	Street  string
}
