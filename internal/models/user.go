package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" binding:"required" gorm:"not null"`
	LastName  string `json:"last_name" binding:"required" gorm:"not null"`
	Email     string `json:"email" binding:"required,email" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"`
	CartId    string `json:"cart_id"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Role      string `json:"role" gorm:"default:user"`
	Age       int    `json:"age"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
}

type UserRegistrationRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=1280"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Age       int    `json:"age"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	CartId    *string `json:"cart_id"`
	Address   *string `json:"address"`
	Type      *string `json:"type"`
	Age       *int    `json:"age"`
	IsActive  *bool   `json:"is_active"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
}

type UserLoginResponse struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}
