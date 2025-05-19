package model

import "time"

type GetUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type User struct {
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ID               string    `gorm:"default:uuid_generate_v4()" json:"id"`
	Name             string    `gorm:"not null" json:"name"`
	Email            string    `gorm:"not null;unique" json:"email"`
	CompanyName      string    `gorm:"not null" json:"company_name"`
	PhoneNumber      string    `gorm:"not null" json:"phone_number"`
	MembershipStatus string    `gorm:"default:PAYG" json:"membership_status"`
}

type CreateUserRequest struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	CompanyName string `json:"company_name" validate:"omitempty"`
	Phone       string `json:"phone_number" validate:"required"`
}

type MembershipStatusChangeRequest struct {
	Status string `json:"membership_status" validate:"required"`
}
