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

// @Summary		Create a new blog post (strictly for admin)
// @Description	Create a new blog post with the provided information
// @Tags			posts
//
// @Security		BearerAuth
//
// @Accept			multipart/form-data
// @Produce		json
// @Param			title		formData	string	true	"Title of the post"
// @Param			slug		formData	string	true	"Slug of the post"
// @Param			content		formData	string	true	"Content of the post"
// @Param			status		formData	string	true	"Status of the post (draft/published)"
// @Param			cover_image	formData	file	true	"Cover image for the post"
// @Success		201			{object}	model.ResponseHTTP{data=model.PostResponse}
// @Failure		400			{object}	model.ResponseHTTP{}
// @Failure		500			{object}	model.ResponseHTTP{}
// @Router			/api/v1/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid form-data request",
			Data:    nil,
		})
	}

	payload := model.PostRequest{
		Title:      form.Value["title"][0],
		Slug:       form.Value["slug"][0],
		Content:    form.Value["content"][0],
		Status:     form.Value["status"][0],
		CoverImage: form.File["cover_image"][0],
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
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Post with this slug/title exists",
				Data:    nil,
			})
		}

		if strings.Contains(err.Error(), "error uploading image") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
				Success: false,
				Message: err.Error(),
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
		Message: "Successfully saved post",
		Data:    *post,
	})
}

// @Summary		Get all published posts
// @Description	Fetch a paginated list of posts from the database
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			page	query		int	false	"Page number (default is 1)"
// @Param			limit	query		int	false	"Number of posts per page (default is 10)"
// @Success		200		{array}		model.ResponseHTTP{data=model.TotalPostResponse}
// @Failure		500		{object}	model.ResponseHTTP{}
// @Router			/api/v1/posts [get]
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

// @Summary		Get all draft posts (strictly for admin)
// @Description	Fetch a paginated list of posts from the database
// @Tags			posts
//
// @Security		BearerAuth
//
// @Accept			json
// @Produce		json
// @Param			page	query		int	false	"Page number (default is 1)"
// @Param			limit	query		int	false	"Number of posts per page (default is 10)"
// @Success		200		{array}		model.ResponseHTTP{data=model.TotalPostResponse}
// @Failure		500		{object}	model.ResponseHTTP{}
// @Router			/api/v1/posts/drafts [get]
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
		if strings.Contains(err.Error(), "post not found") {
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

// @Summary		Update blog post (strictly for admin)
// @Description	Update a blog post with the provided information
// @Tags			posts
//
// @Security		BearerAuth
//
// @Accept			multipart/form-data
// @Produce		json
// @Param			title		formData	string	true	"Title of the post"
// @Param			slug		formData	string	true	"Slug of the post"
// @Param			content		formData	string	true	"Content of the post"
// @Param			status		formData	string	true	"Status of the post (draft/published)"
// @Param			cover_image	formData	file	true	"Cover image for the post"
// @Success		200			{object}	model.ResponseHTTP{data=model.PostResponse}
// @Failure		400			{object}	model.ResponseHTTP{}
// @Failure		404			{object}	model.ResponseHTTP{}
// @Failure		500			{object}	model.ResponseHTTP{}
// @Router			/api/v1/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "post ID is required",
			Data:    nil,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid form-data request",
			Data:    nil,
		})
	}

	payload := model.PostRequest{
		Title:      form.Value["title"][0],
		Slug:       form.Value["slug"][0],
		Content:    form.Value["content"][0],
		Status:     form.Value["status"][0],
		CoverImage: form.File["cover_image"][0],
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	post, err := h.postService.UpdatePost(id, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "post not found") {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Post not found",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully updated post",
		Data:    *post,
	})
}

// UploadImage uploads an image for the post body
//
//	@Summary		Uploads an image for the post body (strictly for admin)
//	@Description	Uploads an image for the post body
//	@Tags			posts
//
//	@Security		BearerAuth
//
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			post_id	formData	string	true	"ID of the post"
//	@Param			image	formData	file	true	"Image to upload"
//	@Success		201		{object}	model.ResponseHTTP{data=model.UploadImageResponse}
//	@Failure		400		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/posts/upload-image [post]
func (h *PostHandler) UploadImage(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid form-data request",
			Data:    nil,
		})
	}

	payload := model.UploadImageRequest{
		PostID: form.Value["post_id"][0],
		Image:  form.File["image"][0],
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	image, err := h.postService.UploadImageFile(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "error uploading image") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
				Success: false,
				Message: err.Error(),
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
		Message: "Successfully uploaded image",
		Data:    *image,
	})
}

// DeletePost deletes a post by ID
//
//	@Summary		Delete a post
//	@Description	Delete a post by ID
//	@Tags			posts
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		204	{object}	model.ResponseHTTP{}
//	@Failure		400	{object}	model.ResponseHTTP{}
//	@Failure		404	{object}	model.ResponseHTTP{}
//	@Failure		500	{object}	model.ResponseHTTP{}
//	@Router			/api/v1/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post ID is required",
		})
	}

	err := h.postService.DeletePost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Post not found",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully deleted post",
		Data:    nil,
	})
}
