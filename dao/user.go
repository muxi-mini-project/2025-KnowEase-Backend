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
	var user models.User
	var usermessage models.Usermessage
	re := ud.db.Where("id = ?", UserID).First(&user)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	usermessage.UserID = UserID
	usermessage.ImageURL = user.ImageURL
	usermessage.Username = user.Username
	return &usermessage, nil
}

// 创建新用户
func (ud *UserDao) CreateNewUser(user *models.User) error {
	re := ud.db.Create(user)
	return re.Error
}

// 修改密码-忘记密码
func (ud *UserDao) ChangePassword(user *models.User) error {
	return ud.db.Save(&user).Error
}

// 查询用户id
func (ud *UserDao) SearchUserid(UserID string) error {
	var User models.User
	return ud.db.Where("id = ?", UserID).First(&User).Error
}

// 修改用户个人主页背景
func (ud *UserDao) ChangeUserBackground(UserID, Newbackground string) error {
	err := ud.db.Model(&models.User{}).Where("id = ?", UserID).Update("page_background_url", Newbackground).Error
	return err
}

// 修改个人头像
func (ud *UserDao) ChangeUserPicture(UserID, NewPicture string) error {
	err := ud.db.Model(&models.User{}).Where("id = ?", UserID).Update("image_url", NewPicture).Error
	return err
}

// 修改密码-个人主页
func (ud *UserDao) ChangeUserPassword(UserID, NewPassword string) error {
	err := ud.db.Model(&models.User{}).Where("id = ?", UserID).Update("password", NewPassword).Error
	return err
}

// 修改用户邮箱
func (ud *UserDao) ChangeUserEmail(UserID, NewEmail string) error {
	err := ud.db.Model(&models.User{}).Where("id = ?", UserID).Update("email", NewEmail).Error
	return err
}

// 修改用户名
func (ud *UserDao) ChangeUsername(UserID, NewName string) error {
	err := ud.db.Model(&models.User{}).Where("id = ?", UserID).Update("username", NewName).Error
	return err
}

// 查找所有用户
func (ud *UserDao) SearchAllUser() ([]string, error) {
	var UserIDs []string
	err := ud.db.Model(&models.User{}).Select("id").Find(&UserIDs).Error
	if err != nil {
		return nil, err
	}
	return UserIDs, nil
}
