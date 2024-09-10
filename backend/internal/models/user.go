package models

import (
	"time"
)

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"type:varchar(100);unique;not null"`
	Email        string    `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	DeletedAt    time.Time `gorm:"index"`
}
