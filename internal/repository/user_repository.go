package repository

import (
	"errors"
	"fmt"

	"github.com/mannanpatel/Restro/common/database"
	"github.com/mannanpatel/Restro/internal/models"
)

type userRepository struct{}

func NewUserRepository() IUserRepository {
	return &userRepository{}
}

// RegisterUser registers a new user
func (ur *userRepository) RegisterUserWithProfile(user *models.User) error {
	// Check if user already exists
	var existing models.User
	err := database.DB.Where("email = ?", user.Email).First(&existing).Error
	if err == nil {
		return errors.New("user already exists")
	}

	// Create user
	result := database.DB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to register user: %w", result.Error)
	}

	return nil
}

// FindUserByEmail finds a user by email
func (ur *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserById fetches a user by ID
func (ur *userRepository) GetUserById(id int) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", id).First(&user).Error
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

	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
