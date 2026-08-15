package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstName" gorm:"not null"`
	LastName  string `json:"lastName" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" gorm:"not null"`
	CartId    string `json:"cartId"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Role      string `json:"role" gorm:"default:'user'"`
	Age       int    `json:"age"`
	IsActive  bool   `json:"isActive" gorm:"default:true"`
}

type UserDetails struct {
	ID        uint   `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty" binding:"required"`
	LastName  string `json:"lastName,omitempty" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password,omitempty" binding:"required,min=6"`
	Role      string `json:"role" binding:"omitempty,oneof=admin user"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
}

type UserLoginResponse struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}
