package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Order struct {
	ID                      string         `gorm:"default:uuid_generate_v4()" json:"id"`
	UserID                  string         `gorm:"not null" json:"user_id"`
	User                    User           `gorm:"foreignKey:UserID" json:"user"`
	ProductName             string         `gorm:"not null" json:"product_name"`
	ProductDescription      string         `gorm:"type:text" json:"product_description"`
	ProductDescriptionImage string         `json:"product_description_image,omitempty"`
	ShootType               string         `gorm:"not null" json:"shoot_type"`
	FinishType              string         `json:"finish_type"`
	Quantity                int            `gorm:"not null" json:"quantity"`
	Details                 datatypes.JSON `gorm:"type:jsonb" json:"details"`
	Shots                   pq.StringArray `gorm:"type:text[]" json:"shots"`
	DeliverySpeed           string         `gorm:"default:Standard" json:"delivery_speed"`
	Status                  string         `gorm:"default:QUOTE RECEIVED" json:"status"`
	CreatedAt               time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type OrderRequest struct {
	UserEmail          string                 `json:"user_email" validate:"required,email"`
	ProductName        string                 `json:"product_name" validate:"required"`
	ProductDescription string                 `json:"product_description" validate:"required"`
	ShootType          string                 `json:"shoot_type" validate:"required"`
	Details            map[string]interface{} `json:"details" validate:"omitempty"`
	FinishType         string                 `json:"finish_type" validate:"omitempty"`
	Quantity           int                    `json:"quantity" validate:"omitempty"`
	Shots              []string               `json:"shots" validate:"omitempty"`
	DeliverySpeed      string                 `json:"delivery_speed" validate:"omitempty"`
	Status             string                 `json:"status" validate:"omitempty"`
}

type OrderStatusChangeRequest struct {
	Status string `json:"status" validate:"required"`
}
