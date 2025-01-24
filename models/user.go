package models

const (
	UserTableName = "user"
)

type User struct {
	ID       string `gorm:"primaryKey;column:id"`
	Password string `gorm:"column:password" json:"password"`
	Email    string `gorm:"column:email" json:"email"`
	Username string `gorm:"column:username"`
	Role     string `gorm:"column:role"`
	ImageURL string `gorm:"type:varchar(255)" json:"url"`
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

func (u *User) TableName() string {
	return UserTableName
}

