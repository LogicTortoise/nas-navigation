package models

import (
	"time"
)

type Link struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CategoryID    uint      `json:"category_id" gorm:"not null"`
	Name          string    `json:"name" gorm:"size:100;not null"`
	URL           string    `json:"url" gorm:"size:500;not null"`
	Icon          string    `json:"icon" gorm:"size:500"`
	Description   string    `json:"description" gorm:"size:500"`
	SortOrder     int       `json:"sort_order" gorm:"default:0"`
	RestartScript string    `json:"restart_script" gorm:"size:500"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Category      *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type LinkRequest struct {
	CategoryID    uint   `json:"category_id" binding:"required"`
	Name          string `json:"name" binding:"required,max=100"`
	URL           string `json:"url" binding:"required,max=500"`
	Icon          string `json:"icon" binding:"max=500"`
	Description   string `json:"description" binding:"max=500"`
	SortOrder     int    `json:"sort_order"`
	RestartScript string `json:"restart_script" binding:"max=500"`
}
