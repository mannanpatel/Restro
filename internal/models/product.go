package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image"`
	IsAvailable bool    `json:"is_available" gorm:"default:true"`
	Rating      float64 `json:"rating" gorm:"type:decimal(3,2);default:0"`
	Type        string  `json:"type"`
	MealTime    string  `json:"meal_time"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Category    *string  `json:"category"`
	ImageURL    *string  `json:"image"`
	IsAvailable *bool    `json:"is_available"`
	Rating      *float64 `json:"rating"`
	Type        *string  `json:"type"`
	MealTime    *string  `json:"meal_time"`
}
