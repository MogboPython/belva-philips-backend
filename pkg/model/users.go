package model

import "time"

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type ContactUsRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}

type User struct {
	ID                string `gorm:"default:uuid_generate_v4()" json:"id"`
	Name              string `gorm:"not null" json:"name"`
	Email             string `gorm:"not null;unique" json:"email"`
	CompanyName       string `gorm:"not null" json:"company_name"`
	Phone             string `gorm:"not null" json:"phone_number"`
	PreferredMode     string `gorm:"not null" json:"preferred_mode_of_communication"`
	WantToReceiveText bool   `gorm:"default:FALSE" json:"want_to_receive_text"`
	// TODO: Add a field for the user's membership status
	MembershipStatus string    `gorm:"default:PAYG" json:"membership_status"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateUserRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	CompanyName       string `json:"company_name"`
	Phone             string `json:"phone_number"`
	PreferredMode     string `json:"preferred_mode_of_communication"`
	WantToReceiveText bool   `json:"want_to_receive_text"`
}
