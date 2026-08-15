package interfaces

import (
	"github.com/mannanpatel/Restro/internal/models"
)

type IUserService interface {
	Login(req models.UserLoginRequest) (*models.UserLoginResponse, error)
	RegisterUser(user *models.UserDetails) error
	// RegisterUser(req models.UserRegisterRequest) error
}
