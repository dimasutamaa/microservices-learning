package models

import "time"

type Wallet struct {
	ID        string    `gorm:"type:varchar(255);primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null;unique"`
	Balance   int64     `gorm:"type:bigint;not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
