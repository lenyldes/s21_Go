package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

var Roles = map[string]bool {
	"admin": true,
}

type User struct {
	Id           uuid.UUID
	Login        string
	Name         string
	Surname      string
	Phonenumber  string
	PasswordHash []byte
	Email        string
	Role         string
}

func validate(login, name, surname, phonenumber, password, email, role string) error {
	if login == "" {
		return ErrInvalidUserData{"Empty login"}
	}
	if name == "" {
		return ErrInvalidUserData{"Empty name"}
	}
	if surname == "" {
		return ErrInvalidUserData{"Empty surname"}
	}
	if (len(password) < 8) {
		return ErrInvalidUserData{"The password must be at least 8 characters long"}
	}
	if _, ok := Roles[role]; !ok {
		return ErrInvalidUserData{"Unexistent role"}
	}

	ok, err := regexp.Match(`^\+7[98]\d{9}$`, []byte(phonenumber))
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidUserData{"Incorrect phone number format"}
	}
	ok, err = regexp.Match(`^\S+@\S+\.\S+$`, []byte(email)) 
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidUserData{"Incorrect email"}
	}

	return nil
}

func NewUser(login, name, surname, phonenumber, password, email, role string) (*User, error) {
	if err := validate(login, name, surname, phonenumber, password, email, role); err != nil {
		return nil, err
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:           uuid.New(),
		Login:        login,
		Name: name,
		Surname: surname,
		Phonenumber: phonenumber,
		PasswordHash: passwordHash,
		Email:        email,
		Role:         role,
	}, nil
}

func hashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (user *User) ChangePassword(newPassword string) error {
	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	return nil
}

type AccessClaims struct {
	TokenType string
	UserId    string
	Role      string
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	TokenType string
	jwt.RegisteredClaims
}
