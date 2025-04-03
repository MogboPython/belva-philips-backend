package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
)

// UserService interface defines methods for user business logic
type UserService interface {
	CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error)
	GetUserByID(id string) (*model.UserResponse, error)
	GetAllUsers(page, limit string) ([]*model.UserResponse, error)
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

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers(pageStr, limitStr string) ([]*model.UserResponse, error) {
	// Convert to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	users, err := s.userRepo.GetAll(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var userResponses []*model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, mapUserToResponse(user))
	}

	return userResponses, nil
}

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
