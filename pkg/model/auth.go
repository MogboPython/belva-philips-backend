package model

type TokenRequestPayload struct {
	UserSessionID string `json:"sessionId"`
}

// AdminLoginRequest represents the request payload for admin login
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
