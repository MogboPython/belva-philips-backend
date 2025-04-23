package service

import (
	"errors"
	"fmt"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"

	"github.com/gofiber/fiber/v2/log"
)

type AdminService interface {
	Login(req *model.AdminLoginRequest) (string, error)
	GetUserByID(id string) (*model.UserResponse, error)
	GetAllUsers(page, limit string) ([]*model.UserResponse, error)
}

type adminService struct {
	userRepo repository.UserRepository
	// postRepo repository.PostRepository
}

func NewAdminService(userRepo repository.UserRepository) AdminService {
	return &adminService{
		userRepo: userRepo,
		// postRepo: postRepo,
	}
}

// Login Admin with username and password
func (*adminService) Login(req *model.AdminLoginRequest) (string, error) {
	if (req.Username != config.Config("ADMIN_USERNAME_HASH")) || (req.Password != config.Config("ADMIN_PASSWORD_HASH")) {
		return "", errors.New("incorrect username or password")
	}
	// FIXME: fix password ish
	// if !utils.CheckPasswordHash(req.Username, config.Config("ADMIN_USERNAME_HASH")) || !utils.CheckPasswordHash(req.Password, config.Config("ADMIN_PASSWORD_HASH")) {
	// 	return "", errors.New("incorrect username or password")
	// }

	token, err := utils.GenerateToken("AdminSession", "admin")
	if err != nil {
		log.Error("Error signing token:", err)
		return "", errors.New("error generating token")
	}

	return token, nil
}

// GetUserByID retrieves a user by ID
func (s *adminService) GetUserByID(id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return mapUserToResponse(user), nil
}

// GetAllUsers retrieves all users
func (s *adminService) GetAllUsers(pageStr, limitStr string) ([]*model.UserResponse, error) {
	// Convert to integers
	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	users, err := s.userRepo.GetAll(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	userResponses := make([]*model.UserResponse, len(users))
	for i, order := range users {
		userResponses[i] = mapUserToResponse(order)
	}

	return userResponses, nil
}
