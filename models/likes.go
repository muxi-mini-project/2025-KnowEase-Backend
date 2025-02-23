package models

import "time"

type UserLikeHistory struct {
	UserID    string `gorm:"not null"`
	PostID    string `gorm:"not null"`
	ID        int    `gorm:"primaryKey"`
	CreatedAt time.Time
}
type UserSaveHistory struct {
	UserID    string `gorm:"not null"`
	PostID    string `gorm:"not null"`
	ID        int    `gorm:"primaryKey"`
	CreatedAt time.Time
}


