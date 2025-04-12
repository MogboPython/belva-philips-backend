package model

import (
	"time"
)

// ResponseHTTP represents response body of this API
type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type UserResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	CompanyName       string    `json:"company_name"`
	Phone             string    `json:"phone_number"`
	PreferredMode     string    `json:"preferred_mode_of_communication"`
	WantToReceiveText bool      `json:"want_to_receive_text"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type OrderResponse struct {
	ID                      string                 `json:"id"`
	UserID                  string                 `json:"user_id"`
	UserEmail               string                 `json:"user_email"`
	UserMembershipStatus    string                 `json:"user_membership_status"`
	ProductName             string                 `json:"product_name"`
	ProductDescription      string                 `gorm:"type:text" json:"product_description"`
	ProductDescriptionImage string                 `json:"product_description_image,omitempty"`
	ShootType               string                 `json:"shoot_type"`
	FinishType              string                 `json:"finish_type"`
	Quantity                int                    `json:"quantity"`
	Details                 map[string]interface{} `json:"details"`
	Shots                   []string               `json:"shots"`
	DeliverySpeed           string                 `json:"delivery_speed"`
	Status                  string                 `json:"status"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
}
