package models

import "time"

type Transaction struct {
	ID          string `gorm:"primaryKey"`
	UserID      string `gorm:"not null;index"`
	MerchantID  string
	Amount      int64  `gorm:"not null"`
	Type        string `gorm:"not null"`
	Status      string `gorm:"not null;index"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
