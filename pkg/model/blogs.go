package model

// import "gorm.io/gorm"

type PostResponse struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Content    string `json:"content"`
	CoverImage string `json:"cover_image"`
	Status     string `json:"status"`
}

// type Post struct {
// 	gorm.Model

// 	Title      string `gorm:"not null"`
// 	Slug       string `gorm:"not null;unique"`
// 	Content    string `gorm:"not null"`
// 	CoverImage string `gorm:"not null"`
// 	Status     string `gorm:"not null"`
// }
