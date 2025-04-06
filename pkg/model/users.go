package model

import "time"

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID                string    `gorm:"default:uuid_generate_v4()" json:"id"`
	Name              string    `gorm:"not null" json:"name"`
	Email             string    `gorm:"not null;unique" json:"email"`
	CompanyName       string    `gorm:"not null" json:"company_name"`
	Phone             string    `gorm:"not null" json:"phone_number"`
	PreferredMode     string    `gorm:"not null" json:"preferred_mode_of_communication"`
	WantToReceiveText bool      `json:"want_to_receive_text"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateUserRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	CompanyName       string `json:"company_name"`
	Phone             string `json:"phone_number"`
	PreferredMode     string `json:"preferred_mode_of_communication"`
	WantToReceiveText bool   `json:"want_to_receive_text"`
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
