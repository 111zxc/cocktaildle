package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStats struct {
	ID                  string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID              string         `gorm:"type:uuid;not null;index"`
	RandomGamesPlayed   int            `gorm:"type:int;not null;default:0"`
	RandomGamesWon      int            `gorm:"type:int;not null;default:0"`
	RandomCurrentStreak int            `gorm:"type:int;not null;default:0"`
	RandomMaxStreak     int            `gorm:"type:int;not null;default:0"`
	DailyGamesPlayed    int            `gorm:"type:int;not null;default:0"`
	DailyGamesWon       int            `gorm:"type:int;not null;default:0"`
	DailyCurrentStreak  int            `gorm:"type:int;not null;default:0"`
	DailyMaxStreak      int            `gorm:"type:int;not null;default:0"`
	LastPlayed          time.Time      `gorm:"type:timestamp"`
	CreatedAt           time.Time      `gorm:"autoCreateTime"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}
