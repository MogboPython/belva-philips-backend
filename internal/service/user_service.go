package service

import (
	"errors"
	"fmt"
	_ "io"
	"log"
	_ "net/http"
	_ "strings"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
)

// UserService interface defines methods for user business logic
type UserService interface {
	CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error)
	GetUserByID(id string) (*model.UserResponse, error)
	GetUserByEmail(req *model.GetUserByEmailRequest) (*model.UserResponse, error)
	// GetAllUsers(page, limit string) ([]*model.UserResponse, error)
	// UpdateUser(id int64, req *model.UpdateUserRequest) (*model.UserResponse, error)
	// DeleteUser(id int64) error
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error) {
	if req == nil {
		return nil, errors.New("invalid request")
	}

	user := &model.User{
		Name:              req.Name,
		Email:             req.Email,
		Phone:             req.Phone,
		CompanyName:       req.CompanyName,
		PreferredMode:     req.PreferredMode,
		WantToReceiveText: req.WantToReceiveText,
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Printf("error saving user: %v", err)
		return nil, err
	}

	return mapUserToResponse(user), nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return mapUserToResponse(user), nil
}

func (s *userService) GetUserByEmail(req *model.GetUserByEmailRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return mapUserToResponse(user), nil
}

// ContactUs sends a contact request to the website admin
// func (s *userService) ContactUs(req *model.ContactUsRequest) error {
// 	url := "https://api.useplunk.com/v1/send"

// 	payload := strings.NewReader("{\n  \"to\": \"<string>\",\n  \"subject\": \"<string>\",\n  \"body\": \"<string>\",\n  \"subscribed\": true,\n  \"name\": \"<string>\",\n  \"from\": \"<string>\",\n  \"reply\": \"<string>\",\n  \"headers\": {}\n}")

// 	email_req, _ := http.NewRequest("POST", url, payload)

// 	email_req.Header.Add("Content-Type", "<content-type>")
// 	email_req.Header.Add("Authorization", "Bearer <token>")

// 	res, err := http.DefaultClient.Do(email_req)
// 	if err != nil {
// 		return fmt.Errorf("failed to find user: %w", err)
// 	}

// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)

// 	fmt.Println(res)
// 	fmt.Println(string(body))
// 	return nil
// }

// UpdateUser updates an existing user
// func (s *userService) UpdateUser(id int64, req *model.UpdateUserRequest) (*model.UserResponse, error) {
// 	if req == nil {
// 		return nil, errors.New("invalid request")
// 	}

// 	existingUser, err := s.userRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Update fields if provided
// 	if req.Username != "" {
// 		existingUser.Username = req.Username
// 	}
// 	if req.Email != "" {
// 		existingUser.Email = req.Email
// 	}
// 	if req.Phone != "" {
// 		existingUser.Phone = req.Phone
// 	}

// 	if err := s.userRepo.Update(id, existingUser); err != nil {
// 		return nil, err
// 	}

// 	return mapUserToResponse(existingUser), nil
// }

// DeleteUser deletes a user by ID
// func (s *userService) DeleteUser(id int64) error {
// 	return s.userRepo.Delete(id)
// }

// mapUserToResponse maps a user model to a user response
func mapUserToResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		Phone:             user.Phone,
		CompanyName:       user.CompanyName,
		PreferredMode:     user.PreferredMode,
		WantToReceiveText: user.WantToReceiveText,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}
