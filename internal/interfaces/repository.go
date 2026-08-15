package interfaces

import (
	"github.com/mannanpatel/Restro/internal/models"
)

type IUserRepository interface {
	FindByEmail(email string) (*models.UserDetails, error)
	// CreateUser(user *models.User) error
}
