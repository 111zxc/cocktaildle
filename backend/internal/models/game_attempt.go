package models

import (
	"time"
)

type GameAttempt struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID        string    `gorm:"type:uuid;not null;index"`
	GameID        string    `gorm:"type:uuid;not null;index"`
	GuessesMade   int       `gorm:"type:int;not null;default:0"`
	Correct       bool      `gorm:"type:boolean;not null;default:false"`
	LastGuessTime time.Time `gorm:"type:timestamp;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
