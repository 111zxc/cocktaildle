package models

import (
	"time"
)

type Guess struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GameAttemptID string    `gorm:"type:uuid;not null;index"`
	Guess         string    `gorm:"type:varchar(255);not null"`
	IsCorrect     bool      `gorm:"type:boolean;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}
