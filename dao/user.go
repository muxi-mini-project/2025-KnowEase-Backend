package dao

import (
	"KnowEase/models"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

// 通过邮箱查找该用户信息
func (ud *UserDao) GetUserFromEmail(email string) (*models.User, error) {
	var user models.User
	re := ud.db.Where("email = ?", email).First(&user)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &user, nil
}

// 通过用户id查找用户信息（用户名，用户头像）
func (ud *UserDao) GetUserFromID(UserID string) (*models.Usermessage, error) {
	var user models.Usermessage
	re := ud.db.Where("id = ?", UserID).First(&user)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &user, nil
}

// 创建新用户
func (ud *UserDao) CreateNewUser(user *models.User) error {
	re := ud.db.Create(user)
	return re.Error
}

// 修改密码
func (ud *UserDao) ChangePassword(user *models.User) error {
	return ud.db.Save(&user).Error
}

// 查询用户id
func (ud *UserDao) SearchUserid(UserID string) error {
	return ud.db.Where(models.User{}, "id = ?", UserID).Error
}
