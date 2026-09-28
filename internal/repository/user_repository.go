package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mannanpatel/Restro/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

// RegisterUser registers a new user
func (ur *userRepository) RegisterUserWithProfile(user *models.User) error {
	// Check if user already exists
	var existing models.User
	err := ur.db.Where("email = ?", user.Email).First(&existing).Error
	if err == nil {
		return errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	// Create user
	result := ur.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to register user: %w", result.Error)
	}

	return nil
}

// FindUserByEmail finds a user by email
func (ur *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserById fetches a user by ID
func (ur *userRepository) GetUserById(id int) (*models.User, error) {
	var user models.User
	err := ur.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (ur *userRepository) UpdateUser(id int, user *models.UpdateUserRequest) error {
	updates := make(map[string]interface{})

	if user.FirstName != nil {
		updates["first_name"] = *user.FirstName
	}

	if user.LastName != nil {
		updates["last_name"] = *user.LastName
	}

	if user.Email != nil {
		updates["email"] = *user.Email
	}

	if user.CartId != nil {
		updates["cart_id"] = *user.CartId
	}

	if user.Address != nil {
		updates["address"] = *user.Address
	}

	if user.Type != nil {
		updates["type"] = *user.Type
	}

	if user.Age != nil {
		updates["age"] = *user.Age
	}

	if user.IsActive != nil {
		updates["is_active"] = *user.IsActive
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := ur.db.Model(&models.User{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (ur *userRepository) DeleteUser(id int) error {
	err := ur.db.Where("id = ?", id).Delete(&models.User{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (ur *userRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := ur.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

func (ur *userRepository) UpdateUserRole(id int, role string) error {
	result := ur.db.Model(&models.User{}).Where("id = ?", id).Update("role", role)
	if result.Error != nil {
		return fmt.Errorf("failed to update user role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
