package models

import "github.com/golang-jwt/jwt"

type Operator struct {
	ID                   int    `json:"id,omitempty"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type RegisterValidationErrors struct {
	NameError                 string `json:"name,omitempty"`
	EmailError                string `json:"email,omitempty"`
	PhoneError                string `json:"phone,omitempty"`
	PasswordError             string `json:"password,omitempty"`
	PasswordConfirmationError string `json:"password_confirmation_error,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
