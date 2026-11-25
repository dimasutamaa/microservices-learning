package models

import "time"

type Wallet struct {
	ID        string    `gorm:"type:varchar(255);primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Balance   int64     `gorm:"type:bigint;not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type WalletTransfer struct {
	TransactionID string `gorm:"primaryKey"`
	FromUserID    string `gorm:"not null;index"`
	ToUserID      string `gorm:"not null;index"`
	Amount        int64  `gorm:"not null"`
	CreatedAt     time.Time
}
