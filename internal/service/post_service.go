package service

import (
	"fmt"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
)

type PostService interface {
	CreatePost(req *model.PostRequest) (*model.PostResponse, error)
	GetPostByID(id string) (*model.PostResponse, error)
	GetAllDrafts(pageStr, limitStr string) ([]*model.PostResponse, error)
	GetAllPosts(page, limit string) ([]*model.PostResponse, error)
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
	coverImageURL, err := uploadImage(req.CoverImage, "blog-cover-photos")

	if err != nil {
		log.Error("error uploading image: %v", err)
		return nil, err
	}

	post := &model.Post{
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

	coverImageURL, err := uploadImage(req.CoverImage, "blog-cover-photos")

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

// mapUserToResponse maps a user model to a user response
func mapPostToResponse(post *model.Post) *model.PostResponse {
	return &model.PostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		CoverImage: PublicImageURL(post.CoverImage),
		Status:     post.Status,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}
