package models

import (
	"time"
)

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Icon      string    `json:"icon" gorm:"size:100"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Links     []Link    `json:"links,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Icon      string `json:"icon" binding:"max=100"`
	SortOrder int    `json:"sort_order"`
}
