package dto

import (
	"fmt"
	"strings"
)

type SignupRequest struct {
	Login       string `json:"login" binding:"required" example:"user"`
	Email       string `json:"email" binding:"required" example:"user@example.com"`
	Role        string `json:"role" binding:"required" example:"admin"`
	Password    string `json:"password" binding:"required" example:"Password123!"`
	Name        string `json:"name" example:"Иван"`
	Surname     string `json:"surname" example:"Иванов"`
	Phonenumber string `json:"phonenumber" example:"+79001234567"`
}

type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (dto SignupRequest) Validate() error {
	if strings.TrimSpace(dto.Login) == "" {
		return fmt.Errorf("поле login обязательно")
	}
	if strings.TrimSpace(dto.Email) == "" {
		return fmt.Errorf("поле email обязательно")
	}
	if strings.TrimSpace(dto.Role) == "" {
		return fmt.Errorf("поле role обязательно")
	}
	if dto.Password == "" {
		return fmt.Errorf("поле password обязательно")
	}
	return nil
}

type SigninRequest struct {
	Login    string `json:"login" binding:"required" example:"user"`
	Password string `json:"password" binding:"required" example:"Password123!"`
}

type SigninResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (dto SigninRequest) Validate() error {
	if strings.TrimSpace(dto.Login) == "" {
		return fmt.Errorf("поле login обязательно")
	}
	if dto.Password == "" {
		return fmt.Errorf("поле password обязательно")
	}
	return nil
}

type ResetPasswordRequest struct {
	Login string `json:"login" binding:"required" example:"user"`
}

func (dto ResetPasswordRequest) Validate() error {
	if strings.TrimSpace(dto.Login) == "" {
		return fmt.Errorf("поле login обязательно")
	}
	return nil
}

type RefreshTokenRequest struct {
	RefreshToken *string `json:"refresh_token" example:"22222222-2222-4222-8222-222222222222"`
}

func (dto RefreshTokenRequest) Validate() error {
	if dto.RefreshToken == nil {
		return fmt.Errorf("поле refresh_token обязательно")
	}
	return nil
}

type RefreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ChangePasswordRequest struct {
	Login           *string `json:"login" example:"user"`
	CurrentPassword *string `json:"current_password" example:"qwertyqq"`
	NewPassword     *string `json:"new_password" example:"asdfghqq"`
}

func (req ChangePasswordRequest) Validate() error {
	if req.Login == nil {
		return fmt.Errorf("field not specified: %s", *req.Login)
	}
	if req.CurrentPassword == nil {
		return fmt.Errorf("field not specified: %s", *req.CurrentPassword)
	}
	if req.NewPassword == nil {
		return fmt.Errorf("field not specified: %s", *req.NewPassword)
	}
	return nil
}
