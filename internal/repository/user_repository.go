package repository

import (
	"github.com/mannanpatel/Restro/common/database"
	"github.com/mannanpatel/Restro/internal/interfaces"
	"github.com/mannanpatel/Restro/internal/models"
)

type userRepository struct{}

func NewUserRepository() interfaces.IUserRepository {
	return &userRepository{}
}

func (ur *userRepository) FindByEmail(email string) (*models.UserDetails, error) {
	var user models.UserDetails
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
