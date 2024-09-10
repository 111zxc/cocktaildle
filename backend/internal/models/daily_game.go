package models

import (
	"time"
)

type DailyGame struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CocktailID    string    `gorm:"type:varchar(100);not null"`
	Date          time.Time `gorm:"type:date;not null;unique"`
	PlayersPlayed int       `gorm:"type:int;not null;default:0"`
	PlayersWon    int       `gorm:"type:int;not null;default:0"`
	AvgGuesses    float64   `gorm:"type:float;not null;default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	DeletedAt     time.Time `gorm:"index"`
}
