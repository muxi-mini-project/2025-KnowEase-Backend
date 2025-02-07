package models

const (
	UserTableName = "user"
)

type User struct {
	ID                string `gorm:"primaryKey;column:id"`
	Password          string `gorm:"column:password" json:"password"`
	Email             string `gorm:"column:email" json:"email"`
	Username          string `gorm:"column:username"`
	Role              string `gorm:"column:role"`
	LikeCount         int    `gorm:"default:0"`
	FollowerCount     int    `gorm:"default:0"`
	FolloweeCount     int    `gorm:"default:0"`
	ImageURL          string `gorm:"default:'https://mini-project.muxixyz.com/lmAp5ighezmJ-vJ7SMnSmCW9Mfau'" json:"url"`
	PageBackgroundURL string `gorm:"default:'https://mini-project.muxixyz.com/lmAp5ighezmJ-vJ7SMnSmCW9Mfau'" json:"backgroundURL"`
}

type Usermessage struct {
	Username string `gorm:"column:username"`
	UserID   string `gorm:"column:id"`
	ImageURL string `gorm:"column:imageurl;type:varchar(255)"`
}
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FollowMessage struct {
	FollowerID string
	FolloweeID string
}


func (u *User) TableName() string {
	return UserTableName
}
