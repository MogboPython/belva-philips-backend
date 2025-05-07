package repository

import (
	"github.com/MogboPython/belvaphilips_backend/internal/storage"

	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	GetByID(postID string) (*model.Post, error)
	GetAllDrafts(offset, limit int) ([]*model.Post, int64, error)
	UpdatePost(post *model.Post) error
	GetAll(offset, limit int) ([]*model.Post, int64, error)
	Delete(postID string) error
}

type postRepository struct {
	db             *gorm.DB
	storageService storage.StorageService
}

func NewPostRepository(db *gorm.DB, storageService storage.StorageService) PostRepository {
	return &postRepository{
		db:             db,
		storageService: storageService,
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
func (r *postRepository) GetAll(offset, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post

	var count int64

	if err := r.db.Where("status = ?", "published").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// GetAllDrafts retrieves all drafts
func (r *postRepository) GetAllDrafts(offset, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post

	var count int64

	if err := r.db.Where("status = ?", "draft").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (r *postRepository) UpdatePost(post *model.Post) error {
	return r.db.Save(&post).Error
}

func (r *postRepository) Delete(postID string) error {
	// Use a transaction to ensure both the database record and file operations are atomic
	return r.db.Transaction(func(tx *gorm.DB) error {
		var post model.Post
		if err := tx.Where("id = ?", postID).First(&post).Error; err != nil {
			return err
		}

		// Delete the post record
		// TODO: also delete folder with post ID
		if err := tx.Delete(&post).Error; err != nil {
			return err
		}

		// If post has a cover image, remove the file
		if post.CoverImage != "" {
			if err := r.storageService.RemoveFile(post.CoverImage); err != nil {
				log.Warnf("Failed to delete cover image %s for post %s: %v",
					post.CoverImage, postID, err)
			}
		}

		return nil
	})
}
