package model

import "time"

type GetUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type User struct {
	ID                string    `gorm:"default:uuid_generate_v4()" json:"id"`
	Name              string    `gorm:"not null" json:"name"`
	Email             string    `gorm:"not null;unique" json:"email"`
	CompanyName       string    `gorm:"not null" json:"company_name"`
	Phone             string    `gorm:"not null" json:"phone_number"`
	PreferredMode     string    `gorm:"not null" json:"preferred_mode_of_communication"`
	WantToReceiveText bool      `gorm:"default:FALSE" json:"want_to_receive_text"`
	MembershipStatus  string    `gorm:"default:PAYG" json:"membership_status"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateUserRequest struct {
	Name              string `json:"name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	CompanyName       string `json:"company_name" validate:"omitempty"`
	Phone             string `json:"phone_number" validate:"required"`
	PreferredMode     string `json:"preferred_mode_of_communication" validate:"omitempty"`
	WantToReceiveText bool   `json:"want_to_receive_text" validate:"omitempty"`
}

type MembershipStatusChangeRequest struct {
	Status string `json:"membership_status" validate:"required"`
}
