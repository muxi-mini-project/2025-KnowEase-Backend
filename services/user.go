package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserDao *dao.UserDao
}

func NewUserService(UserDao *dao.UserDao) *UserService {
	return &UserService{UserDao: UserDao}
}

// 用户注册
func (us *UserService) Register(user *models.User) error {
	encryptedpassword, err := EncryptPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to envrypt password:%w", err)
	}
	user.Password = encryptedpassword
	err = us.UserDao.CreateNewUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user:%w", err)
	}
	return nil
}

// 用户密码登录
func (us *UserService) LoginByPassword(LoginMessqge models.Login) (*models.User, error) {
	User, err := us.UserDao.GetUserFromEmail(LoginMessqge.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find usermessage")
	} else if User == nil {
		return nil, fmt.Errorf("this user is not registered")
	}
	return User, nil
}

// 密码加密
func EncryptPassword(password string) (string, error) {
	encryptedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encryptedpassword), nil
}

// 密码比对
func (us *UserService) ComparePassword(password1, password2 string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	return err
}

// 修改密码
func (us *UserService) ChangePassword(user models.Login) error {
	User, err := us.UserDao.GetUserFromEmail(user.Email)
	if err != nil {
		return fmt.Errorf("failed to find usermessage")
	}
	(*User).Password, err = EncryptPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password")
	}
	if err := us.UserDao.ChangePassword(User); err != nil {
		return fmt.Errorf("failed to change password")
	}
	return nil
}
