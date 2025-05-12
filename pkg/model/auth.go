package model

type TokenRequestPayload struct {
	UserSessionID string `json:"sessionId"`
}

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
