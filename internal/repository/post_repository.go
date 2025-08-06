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

	CreateGallery(gallery *model.Gallery) error
	GetAllGalleries(offset, limit int) ([]*model.Gallery, int64, error)
	GetGalleryByID(id string) (*model.Gallery, error)
	GetGalleryBySlug(slug string) (*model.Gallery, error)
	UpdateGallery(gallery *model.Gallery) error
	DeleteGallery(id string) error
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

		if err := tx.Delete(&post).Error; err != nil {
			return err
		}

		if post.CoverImage != "" {
			if err := r.storageService.RemoveFolder("blog-cover-photos", post.ID); err != nil {
				log.Warnf("Failed to delete cover image %s for post %s: %v", post.CoverImage, postID, err)
			}
		}

		if err := r.storageService.RemoveFolder("blog-body-photos", post.Slug); err != nil {
			log.Warnf("Failed to delete body image for post %s: %v", postID, err)
		}

		return nil
	})
}

func (r *postRepository) CreateGallery(gallery *model.Gallery) error {
	err := r.db.Create(&gallery).Error
	return err
}

func (r *postRepository) GetGalleryByID(id string) (*model.Gallery, error) {
	var gallery model.Gallery

	err := r.db.Where("id = ?", id).First(&gallery).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gallery not found")
		}

		return nil, err
	}

	return &gallery, nil
}

func (r *postRepository) GetGalleryBySlug(slug string) (*model.Gallery, error) {
	var gallery model.Gallery

	err := r.db.Where("slug = ?", slug).First(&gallery).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gallery object not found")
		}

		return nil, err
	}

	return &gallery, nil
}

func (r *postRepository) GetAllGalleries(offset, limit int) ([]*model.Gallery, int64, error) {
	var galleries []*model.Gallery

	var count int64

	if err := r.db.Offset(offset).Limit(limit).Find(&galleries).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return galleries, count, nil
}

func (r *postRepository) UpdateGallery(gallery *model.Gallery) error {
	err := r.db.Save(gallery).Error
	if err != nil {
		if isDuplicateError(err) {
			return errors.New("slug already exists")
		}

		return err
	}

	return nil
}

func (r *postRepository) DeleteGallery(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var gallery model.Gallery
		if err := tx.Where("id = ?", id).First(&gallery).Error; err != nil {
			return err
		}

		return tx.Delete(&gallery).Error
	})
}
