package model

import (
	"mime/multipart"
	"time"
)

type PostRequest struct {
	CoverImage *multipart.FileHeader `form:"cover_image" validate:"omitempty"`
	Title      string                `form:"title" validate:"required"`
	Slug       string                `form:"slug" validate:"required"`
	Content    string                `form:"content" validate:"omitempty"`
	Status     string                `form:"status" validate:"omitempty"`
}

type Post struct {
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ID         string    `gorm:"default:uuid_generate_v4()" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Slug       string    `gorm:"unique" json:"slug"`
	Content    string    `json:"content"`
	CoverImage string    `json:"cover_image"`
	Status     string    `gorm:"default:draft" json:"status"`
}

type UploadImageRequest struct {
	Image  *multipart.FileHeader `form:"image" validate:"required"`
	PostID string                `form:"post_id" validate:"required"`
}

type DeleteImageRequest struct {
	FileName string `json:"file_name"`
}
