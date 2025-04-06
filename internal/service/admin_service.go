package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AdminService interface {
	Login(req *model.AdminLoginRequest) (string, error)
	GetUserByID(id string) (*model.UserResponse, error)
	GetAllUsers(page, limit string) ([]*model.UserResponse, error)
}

type adminService struct {
	userRepo repository.UserRepository
}

func NewAdminService(userRepo repository.UserRepository) AdminService {
	return &adminService{
		userRepo: userRepo,
	}
}

// Login Admin with username and password
func (s *adminService) Login(req *model.AdminLoginRequest) (string, error) {

	if !utils.CheckPasswordHash(req.Username, config.Config("ADMIN_USERNAME_HASH")) || !utils.CheckPasswordHash(req.Password, config.Config("ADMIN_PASSWORD_HASH")) {
		return "", errors.New("incorrect username or password")
	}

	claims := jwt.MapClaims{
		"sessionId": "AdminSession",
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token, sign and generate encoded token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config("ADMIN_JWT_SECRET")))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", errors.New("error generating token")
	}

	return t, nil
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
