package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/internal/storage"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type PostService interface {
	CreatePost(req *model.PostRequest) (*model.PostResponse, error)
	GetPostByID(id string) (*model.PostResponse, error)
	GetAllDrafts(pageStr, limitStr string) (model.TotalPostResponse, error)
	GetAllPosts(page, limit string) (model.TotalPostResponse, error)
	UpdatePost(id string, update *model.PostUpdateRequest) (*model.PostResponse, error)
	UploadImageFile(req *model.UploadImageRequest) (*model.UploadImageResponse, error)
	DeletePost(id string) error
}

type postService struct {
	postRepo       repository.PostRepository
	storageService storage.StorageService
}

func NewPostService(postRepo repository.PostRepository, storageService storage.StorageService) PostService {
	return &postService{
		postRepo:       postRepo,
		storageService: storageService,
	}
}

func (s *postService) CreatePost(req *model.PostRequest) (*model.PostResponse, error) {
	postID := uuid.New()

	coverImageURL, err := s.storageService.UploadFile(req.CoverImage, "blog-cover-photos", postID.String())

	if err != nil {
		return nil, fmt.Errorf("error uploading image: %w", err)
	}

	post := &model.Post{
		ID:         postID.String(),
		Title:      req.Title,
		CoverImage: coverImageURL,
		Slug:       req.Slug,
		Content:    req.Content,
		Status:     req.Status,
	}

	if err := s.postRepo.Create(post); err != nil {
		log.Error("error saving post: ", err)
		return nil, err
	}

	return mapPostToResponse(post), nil
}

func (s *postService) GetAllPosts(pageStr, limitStr string) (model.TotalPostResponse, error) {
	var totalPostResponse model.TotalPostResponse

	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	posts, count, err := s.postRepo.GetAll(offset, limit)
	if err != nil {
		return totalPostResponse, fmt.Errorf("failed to get posts: %w", err)
	}

	postResponses := make([]*model.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = mapPostToResponse(post)
	}

	totalPostResponse.Posts = postResponses
	totalPostResponse.Total = count

	return totalPostResponse, nil
}

func (s *postService) GetAllDrafts(pageStr, limitStr string) (model.TotalPostResponse, error) {
	var totalPostResponse model.TotalPostResponse

	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	posts, count, err := s.postRepo.GetAllDrafts(offset, limit)
	if err != nil {
		return totalPostResponse, fmt.Errorf("failed to get posts: %w", err)
	}

	postResponses := make([]*model.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = mapPostToResponse(post)
	}

	totalPostResponse.Posts = postResponses
	totalPostResponse.Total = count

	return totalPostResponse, nil
}

func (s *postService) GetPostByID(id string) (*model.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		log.Error("failed to find post:", err)
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	return mapPostToResponse(post), nil
}

func (s *postService) UpdatePost(id string, update *model.PostUpdateRequest) (*model.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	if update.CoverImage != nil {
		newCoverImageURL, err := s.storageService.UploadFile(update.CoverImage, "blog-cover-photos", post.ID)
		if err != nil {
			log.Errorf("Failed to upload new cover image: %v", err)
			return nil, fmt.Errorf("failed to upload cover image: %w", err)
		}

		if post.CoverImage != "" {
			if err := s.storageService.RemoveFile(post.CoverImage); err != nil {
				log.Warnf("Failed to delete old cover image %s for post %s: %v", post.CoverImage, post.ID, err)
			}
		}

		post.CoverImage = newCoverImageURL
	}

	post.Title = update.Title
	post.Slug = update.Slug
	post.Content = update.Content
	post.Status = update.Status
	post.UpdatedAt = time.Now()

	if err := s.postRepo.Update(post); err != nil {
		log.Error("error saving post: ", err)
		return nil, err
	}

	return mapPostToResponse(post), nil
}

func (s *postService) UploadImageFile(req *model.UploadImageRequest) (*model.UploadImageResponse, error) {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[req.Image.Header.Get("Content-Type")] {
		return nil, errors.New("error uploading image: Invalid file type. Only JPEG, PNG, and GIF are allowed")
	}

	fileName, err := s.storageService.UploadFile(req.Image, "blog-body-photos", req.PostID)
	if err != nil {
		log.Error("error uploading image: %v", err)
		return nil, fmt.Errorf("error uploading image: %w", err)
	}

	publicURL := utils.PublicImageURL(fileName)

	return &model.UploadImageResponse{
		ImageURL: publicURL,
		FileName: fileName,
	}, nil
}

// func (*postService) DeleteImageFile(req *model.DeleteImageRequest) error {
// 	if err := removeFile(req.FileName); err != nil {
// 		log.Error("error deleting image: %v", err)
// 		return errors.New("error deleting image")
// 	}

// 	return nil
// }

// DeletePost deletes a post
func (s *postService) DeletePost(id string) error {
	err := s.postRepo.Delete(id)
	if err != nil {
		log.Errorf("Failed to delete post %s: %v", id, err)
		return err
	}

	return nil
}

func mapPostToResponse(post *model.Post) *model.PostResponse {
	return &model.PostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		CoverImage: utils.PublicImageURL(post.CoverImage),
		Status:     post.Status,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}
