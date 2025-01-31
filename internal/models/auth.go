package models

import "github.com/golang-jwt/jwt"

type Operator struct {
	ID                   int           `json:"id,omitempty"`
	Role                 Role          `json:"role"`
	Name                 string        `json:"name"`
	UserName             string        `json:"username"`
	Email                *string       `json:"email"`
	Phone                string        `json:"phone"`
	IsActive             bool          `json:"is_active"`
	IsLimitedEventAccess bool          `json:"is_limited_event_access"`
	IsLimitedClassAccess bool          `json:"is_limited_classes_access"`
	DeviceId             *string       `json:"device_id"`
	DeviceAccessToken    *string       `json:"device_access_token"`
	Password             string        `json:"password"`
	PasswordConfirmation string        `json:"password_confirmation,omitempty"`
	AccessEvents         []AccessEvent `json:"access_events,omitempty"`
}

type EventClass struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

type AccessEvent struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	Classes []EventClass `json:"classes"`
}

type Role struct {
	ID   *int    `json:"id"`
	Name *string `json:"name"`
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

type LoginResponseData struct {
	Token string   `json:"token"`
	User  Operator `json:"user"`
}
