package models

import (
	"time"
)

type Script struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Command     string    `json:"command" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Timeout     int       `json:"timeout" gorm:"default:30"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ScriptRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Command     string `json:"command" binding:"required"`
	Description string `json:"description" binding:"max=500"`
	Timeout     int    `json:"timeout"`
}

type ScriptExecuteResult struct {
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
	Duration int64  `json:"duration_ms"`
}
