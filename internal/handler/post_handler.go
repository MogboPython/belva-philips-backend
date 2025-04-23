package handler

import (
	"errors"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PostHandler struct {
	postService service.PostService
	validator   *validator.Validator
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
		validator:   validator.New(),
	}
}

// CreatePost creates a new post
//
//	@Summary		Create a new blog post (strictly for admin)
//	@Description	Create a new blog post with the provided information
//	@Tags			posts
//
//	@Security		BearerAuth
//
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			title		formData	string	true	"Title of the post"
//	@Param			slug		formData	string	true	"Slug of the post"
//	@Param			content		formData	string	true	"Content of the post"
//	@Param			status		formData	string	true	"Status of the post (draft/published)"
//	@Param			cover_image	formData	file	true	"Cover image for the post"
//	@Success		201			{object}	model.ResponseHTTP{data=model.PostResponse}
//	@Failure		400			{object}	model.ResponseHTTP{}
//	@Failure		500			{object}	model.ResponseHTTP{}
//	@Router			/api/v1/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var payload model.PostRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	post, err := h.postService.CreatePost(&payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully saved post",
		Data:    *&post,
	})
}

// GetAllPosts is a function to get all post data from the database
//
//	@Summary		Get all published posts
//	@Description	Fetch a paginated list of posts from the database
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number (default is 1)"
//	@Param			limit	query		int	false	"Number of posts per page (default is 10)"
//	@Success		200		{array}		model.ResponseHTTP{data=[]model.PostResponse}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	posts, err := h.postService.GetAllPosts(pageStr, limitStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully retrieved posts.",
		Data:    posts,
	})
}

// GetAllDraftPosts is a function to get all draft post data from the database
//
//	@Summary		Get all draft posts (strictly for admin)
//	@Description	Fetch a paginated list of posts from the database
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number (default is 1)"
//	@Param			limit	query		int	false	"Number of posts per page (default is 10)"
//	@Success		200		{array}		model.ResponseHTTP{data=[]model.PostResponse}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/posts/drafts [get]
func (h *PostHandler) GetAllDraftPosts(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	posts, err := h.postService.GetAllDrafts(pageStr, limitStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully retrieved posts.",
		Data:    posts,
	})
}

// GetPostByID is a function to get an post by ID
//
//	@Summary		Get post by ID
//	@Description	Get post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	model.ResponseHTTP{data=model.PostResponse}
//	@Failure		404	{object}	model.ResponseHTTP{}
//	@Failure		500	{object}	model.ResponseHTTP{}
//	@Router			/api/v1/posts/{id} [get]
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Post not found",
				Data:    nil,
			})
		}

		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Invalid route",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully found post",
		Data:    *post,
	})
}

// validStatuses := map[string]bool{"draft": true, "published": true}
// if !validStatuses[req] {
// 	return nil, errors.New("invalid request")
// }
