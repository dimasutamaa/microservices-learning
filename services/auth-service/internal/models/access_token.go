package models

import "time"

type AccessToken struct {
	ID        string    `gorm:"type:varchar(255);primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null;index"`
	Revoked   bool      `gorm:"nul null;default:false"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
