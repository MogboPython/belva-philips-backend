package repository

import (
	"errors"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/storage"
	"github.com/lib/pq"

	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	GetByID(postID string) (*model.Post, error)
	GetAllDrafts(offset, limit int) ([]*model.Post, int64, error)
	Update(post *model.Post) error
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

func (r *postRepository) Create(post *model.Post) error {
	err := r.db.Create(&post).Error
	return err
}

func (r *postRepository) GetByID(id string) (*model.Post, error) {
	var post model.Post

	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}

		return nil, err
	}

	return &post, nil
}

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

func (r *postRepository) Update(post *model.Post) error {
	err := r.db.Save(post).Error
	if err != nil {
		if isDuplicateError(err) {
			return errors.New("slug already exists")
		}

		return err
	}

	return nil
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}

	return strings.Contains(err.Error(), "duplicate key")
}

func (r *postRepository) Delete(postID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var post model.Post
		if err := tx.Where("id = ?", postID).First(&post).Error; err != nil {
			return err
		}

		// TODO: also delete folder with post ID
		if err := tx.Delete(&post).Error; err != nil {
			return err
		}

		if post.CoverImage != "" {
			if err := r.storageService.RemoveFile(post.CoverImage); err != nil {
				log.Warnf("Failed to delete cover image %s for post %s: %v",
					post.CoverImage, postID, err)
			}
		}

		return nil
	})
}
