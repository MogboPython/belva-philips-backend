package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type OrdersCount struct {
	Total        int64 `json:"total_orders"`
	ActiveCount  int64 `json:"active_orders"`
	PendingCount int64 `json:"pending_orders"`
}

type Order struct {
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	User               User           `gorm:"foreignKey:UserID" json:"user"`
	ID                 string         `gorm:"default:uuid_generate_v4()" json:"id"`
	UserID             string         `gorm:"type:uuid;not null" json:"user_id"`
	ProductName        string         `gorm:"not null" json:"product_name"`
	ProductDescription string         `gorm:"type:text" json:"product_description"`
	ShootType          string         `gorm:"not null" json:"shoot_type"`
	FinishType         string         `json:"finish_type"`
	DeliverySpeed      string         `gorm:"default:Standard" json:"delivery_speed"`
	Status             string         `gorm:"default:QUOTE RECEIVED" json:"status"`
	Details            datatypes.JSON `gorm:"type:jsonb" json:"details"`
	Shots              pq.StringArray `gorm:"type:text[]" json:"shots"`
	Quantity           int            `gorm:"not null" json:"quantity"`
}

type OrderRequest struct {
	UserID             string         `json:"user_id" validate:"required"`
	ProductName        string         `json:"product_name" validate:"required"`
	ProductDescription string         `json:"product_description" validate:"required"`
	ShootType          string         `json:"shoot_type" validate:"required"`
	FinishType         string         `json:"finish_type" validate:"omitempty"`
	DeliverySpeed      string         `json:"delivery_speed" validate:"omitempty"`
	Status             string         `json:"status" validate:"omitempty"`
	Details            map[string]any `json:"details" validate:"omitempty"`
	Shots              []string       `json:"shots" validate:"omitempty"`
	Quantity           int            `json:"quantity" validate:"omitempty"`
}

type OrderStatusChangeRequest struct {
	Status string `json:"status" validate:"required"`
}
