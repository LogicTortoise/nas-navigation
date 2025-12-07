package models

import (
	"time"
)

type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"size:100;uniqueIndex;not null"`
	Value     string    `json:"value" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SettingRequest struct {
	Key   string `json:"key" binding:"required,max=100"`
	Value string `json:"value"`
}

// 预定义设置键
const (
	SettingUserName    = "user_name"
	SettingGreeting    = "greeting"
	SettingTheme       = "theme"
	SettingLanguage    = "language"
)
