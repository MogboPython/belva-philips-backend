package repository

import (
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	GetByID(postID string) (*model.Post, error)
	GetAllDrafts(offset, limit int) ([]*model.Post, error)
	UpdatePost(post *model.Post) error
	GetAll(offset, limit int) ([]*model.Post, error)
	// Delete(id int64) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

// Create inserts a new post into the database
func (r *postRepository) Create(post *model.Post) error {
	err := r.db.Create(&post).Error
	return err
}

// GetByID retrieves a post by ID
func (r *postRepository) GetByID(id string) (*model.Post, error) {
	var post model.Post

	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// GetAll retrieves all published posts
func (r *postRepository) GetAll(offset, limit int) ([]*model.Post, error) {
	var posts []*model.Post

	if err := r.db.Where("status = ?", "published").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// GetAllDrafts retrieves all drafts
func (r *postRepository) GetAllDrafts(offset, limit int) ([]*model.Post, error) {
	var posts []*model.Post

	if err := r.db.Where("status = ?", "draft").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) UpdatePost(post *model.Post) error {
	return r.db.Save(&post).Error
}
