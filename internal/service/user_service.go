package service

import (
	"fmt"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"

	"github.com/gofiber/fiber/v2/log"
)

// UserService interface defines methods for user business logic
type UserService interface {
	CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error)
	GetUserByID(id string) (*model.UserResponse, error)
	UpdateUserMembershipStatusChange(userID string, request *model.MembershipStatusChangeRequest) (*model.UserResponse, error)
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
	user := &model.User{
		ID:          req.ID,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		CompanyName: req.CompanyName,
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Error("error saving user: %v", err)
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

func (s *userService) UpdateUserMembershipStatusChange(userID string, request *model.MembershipStatusChangeRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.UpdateMembership(userID, request.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return mapUserToResponse(user), nil
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

func mapUserToResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		CompanyName: user.CompanyName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
