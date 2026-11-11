package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:varchar(255);primaryKey"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
