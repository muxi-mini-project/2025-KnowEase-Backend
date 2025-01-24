package models

import "time"

type UserLikes struct {
	UserID        string `gorm:"not null"`
	FollowerCount int    `gorm:"default:0"`
	FolloweeCount int    `gorm:"default:0"`
	LikeCount     int    `gorm:"default:0"`
}
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


