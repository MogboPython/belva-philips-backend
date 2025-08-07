package model

import (
	"time"

	"github.com/lib/pq"
)

type Gallery struct {
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	ID        string         `gorm:"default:uuid_generate_v4()" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Slug      string         `gorm:"unique" json:"slug"`
	Images    pq.StringArray `gorm:"type:text[]" json:"images"`
}

type GalleryRequest struct {
	Title  string   `json:"title" validate:"required"`
	Slug   string   `json:"slug" validate:"required"`
	Images []string `json:"images" validate:"required"`
}

type GalleryUpdateRequest struct {
	Title  string   `json:"title" validate:"required"`
	Slug   string   `json:"slug" validate:"required"`
	Images []string `json:"images" validate:"omitempty"`
}

type GalleryDeleteRequest struct {
	PublicIDs []string `json:"public_urls" validate:"omitempty"`
}

type GalleryResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	ID        string    `json:"id"`
	Images    []string  `json:"images" validate:"omitempty"`
}

type TotalGalleryResponse struct {
	Galleries []*GalleryResponse `json:"galleries"`
	Total     int64              `json:"total"`
}
