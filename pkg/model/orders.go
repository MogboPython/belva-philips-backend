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
	// Props                   string         `json:"props"`
	// Backdrop                string         `json:"backdrop"`
	// ItemsInFrame            string         `json:"items_in_frame"`
	// Shadow                  string         `json:"shadow"`
	// ModelChoice             string         `json:"model_choice"`
	// ModelDisplay            string         `json:"model_display"`
	// VideoType               string         `gorm:"default:Standard" json:"video_type"`
	// VideoQuantity           int            `gorm:"not null" json:"video_quantity"`
	// AnimationPackage        string         `json:"animation_package"`
}

type OrderRequest struct {
	UserEmail          string                 `json:"user_email"`
	ProductName        string                 `json:"product_name"`
	ProductDescription string                 `json:"product_description"`
	ShootType          string                 `json:"shoot_type"`
	Details            map[string]interface{} `json:"details"`
	FinishType         string                 `json:"finish_type,omitempty"`
	Quantity           int                    `json:"quantity,omitempty"`
	Shots              []string               `json:"shots,omitempty"`
	DeliverySpeed      string                 `json:"delivery_speed,omitempty"`
	Status             string                 `json:"status,omitempty"`

	// Details            datatypes.JSON `json:"details,omitempty"`
	// Props        string   `json:"props,omitempty"`
	// Backdrop     string   `json:"backdrop,omitempty"`
	// ItemsInFrame string   `json:"items_in_frame,omitempty"`
	// Shadow       string   `json:"shadow,omitempty"`

	// ModelChoice  string `json:"model_choice,omitempty"`
	// ModelDisplay string `json:"model_display,omitempty"`

	// VideoType        string `json:"video_type,omitempty"`
	// VideoQuantity    int    `json:"video_quantity,omitempty"`
	// AnimationPackage string `json:"animation_package,omitempty"`
}
