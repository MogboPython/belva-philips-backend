package model

import (
	"time"
)

type ResponseHTTP struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type UserResponse struct {
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	CompanyName       string    `json:"company_name"`
	Phone             string    `json:"phone_number"`
	PreferredMode     string    `json:"preferred_mode_of_communication"`
	WantToReceiveText bool      `json:"want_to_receive_text"`
}

type OrderResponse struct {
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	ProductDescription   string         `gorm:"type:text" json:"product_description"`
	ID                   string         `json:"id"`
	UserID               string         `json:"user_id"`
	UserEmail            string         `json:"user_email"`
	UserMembershipStatus string         `json:"user_membership_status"`
	ProductName          string         `json:"product_name"`
	ShootType            string         `json:"shoot_type"`
	FinishType           string         `json:"finish_type"`
	DeliverySpeed        string         `json:"delivery_speed"`
	Status               string         `json:"status"`
	MembershipType       string         `json:"membership_type"`
	Details              map[string]any `json:"details"`
	Shots                []string       `json:"shots"`
	Quantity             int            `json:"quantity"`
}

type TotalOrderResponse struct {
	Orders       []*OrderResponse `json:"orders"`
	OrderNumbers OrdersCount      `json:"orders_count"`
}

type PostResponse struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	CoverImage string    `json:"cover_image"`
	Status     string    `json:"status"`
	ID         string    `json:"id"`
}

type TotalPostResponse struct {
	Posts []*PostResponse `json:"posts"`
	Total int64           `json:"total"`
}

type UploadImageResponse struct {
	ImageURL string `json:"image_url"`
	FileName string `json:"file_name"`
}
