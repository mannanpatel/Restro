package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mannanpatel/Restro/config"
	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository repository.IUserRepository
}

func NewUserService(repository repository.IUserRepository) IUserService {
	return &userService{repository: repository}
}

func (s *userService) Register(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	user.ID = 0
	user.Role = "user"
	user.IsActive = true

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	err = s.repository.RegisterUserWithProfile(user)
	if err != nil {
		return err
	}

	return nil
}

// Login authenticates a user with email and password
func (s *userService) Login(req models.UserLoginRequest) (*models.UserLoginResponse, error) {
	user, err := s.repository.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !user.IsActive {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.AppConfig.Server.JwtSecret))
	if err != nil {
		return nil, err
	}

	return &models.UserLoginResponse{
		Email:  user.Email,
		Role:   user.Role,
		UserId: fmt.Sprintf("%d", user.ID),
		Token:  signed,
	}, nil
}

// GetUserById fetches a user by ID
func (s *userService) GetUserById(id int) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repository.GetUserById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil

}

func (s *userService) UpdateUser(id int, user *models.UpdateUserRequest) error {

	if id == 0 {
		return fmt.Errorf("invalid user id")
	}

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	err := s.repository.UpdateUser(id, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}

	if err := s.repository.DeleteUser(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.User, error) {
	users, err := s.repository.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

func (s *userService) UpdateUserRole(id int, role string) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}
	if role == "" || (role != "admin" && role != "user") {
		return fmt.Errorf("invalid role")
	}

	err := s.repository.UpdateUserRole(id, role)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}
