package model

type TokenRequestPayload struct {
	UserSessionID string `json:"sessionId"`
}

// AdminLoginRequest represents the request payload for admin login
type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
