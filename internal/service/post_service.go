package service

import (
	"errors"
	"fmt"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type PostService interface {
	CreatePost(req *model.PostRequest) (*model.PostResponse, error)
	GetPostByID(id string) (*model.PostResponse, error)
	GetAllDrafts(pageStr, limitStr string) ([]*model.PostResponse, error)
	GetAllPosts(page, limit string) ([]*model.PostResponse, error)
	UploadImageFile(req *model.UploadImageRequest) (*model.UploadImageResponse, error)
	DeletePost(id string) error
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

// CreatePost creates a new post
func (s *postService) CreatePost(req *model.PostRequest) (*model.PostResponse, error) {
	postID := uuid.New()

	coverImageURL, err := uploadFile(req.CoverImage, "blog-cover-photos", postID.String())

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

// GetAllPosts retrieves all published posts
func (s *postService) GetAllPosts(pageStr, limitStr string) ([]*model.PostResponse, error) {
	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	posts, err := s.postRepo.GetAll(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	postResponses := make([]*model.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = mapPostToResponse(post)
	}

	return postResponses, nil
}

// GetAllDrafts retrieves draft posts
func (s *postService) GetAllDrafts(pageStr, limitStr string) ([]*model.PostResponse, error) {
	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	posts, err := s.postRepo.GetAllDrafts(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	postResponses := make([]*model.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = mapPostToResponse(post)
	}

	return postResponses, nil
}

// GetPostByID retrieves an post by ID
func (s *postService) GetPostByID(id string) (*model.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		log.Error("failed to find post:", err)
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	return mapPostToResponse(post), nil
}

// TODO: a lot to work on here
// UpdatePost updates an existing post
func (s *postService) UpdatePost(id string, req *model.PostRequest) (*model.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	// TODO: Check if there was a cover image before then delete now

	coverImageURL, err := uploadFile(req.CoverImage, "blog-cover-photos")

	if err != nil {
		log.Error("error uploading image: %v", err)
		return nil, err
	}

	// Update post fields
	post.Title = req.Title
	post.Slug = req.Slug
	post.Content = req.Content
	post.CoverImage = coverImageURL
	post.Status = req.Status

	return mapPostToResponse(post), nil
}

func (*postService) UploadImageFile(req *model.UploadImageRequest) (*model.UploadImageResponse, error) {
	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[req.Image.Header.Get("Content-Type")] {
		return nil, errors.New("error uploading image: Invalid file type. Only JPEG, PNG, and GIF are allowed")
	}

	// Upload the image to Supabase, let the post ID be the folder name
	fileName, err := uploadFile(req.Image, "blog-body-photos", req.PostID)
	if err != nil {
		log.Error("error uploading image: %v", err)
		return nil, fmt.Errorf("error uploading image: %w", err)
	}

	// Get the public URL of the uploaded image
	publicURL := publicImageURL(fileName)

	return &model.UploadImageResponse{
		ImageURL: publicURL,
		FileName: fileName,
	}, nil
}

// func (*postService) DeleteImageFile(req *model.DeleteImageRequest) error {
// 	// Delete the image to Supabase
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

// mapUserToResponse maps a user model to a user response
func mapPostToResponse(post *model.Post) *model.PostResponse {
	return &model.PostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		CoverImage: publicImageURL(post.CoverImage),
		Status:     post.Status,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}
