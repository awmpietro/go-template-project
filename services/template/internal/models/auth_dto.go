package models

type ExchangeTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LoginRequest sent by mobile with Firebase Token
type LoginRequest struct {
	FirebaseToken string `json:"firebase_token" validate:"required"`
}

// LoginResponse contains the app token + user basic info
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"` // Imported from user_dto.go
}

// ResetPasswordRequest for password reset flow
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// DirectLoginRequest (optional) if allowing backend auth
type DirectLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
