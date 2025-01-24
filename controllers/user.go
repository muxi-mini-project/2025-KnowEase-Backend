package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserControllers struct {
	UserService  *services.UserService
	EmailService *services.EmailService
	TokenService *services.TokenService
}

func NewUserControllers(UserService *services.UserService, EmailService *services.EmailService, TokenService *services.TokenService) *UserControllers {
	return &UserControllers{UserService: UserService, EmailService: EmailService, TokenService: TokenService}
}

// Sendemail 发送验证码邮件-注册
// @Summary 发送验证码
// @Description 接收用户的邮箱地址后判断该用户是否已注册并发送验证码邮件
// @Tags 用户注册
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 200 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 409 {object} models.Response "该用户已注册"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/register/sendemail [post]
func (u *UserControllers) RegisterSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.Register(email.Email); err != nil {
		c.JSON(http.StatusConflict, models.Write(err.Error()))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		response := models.Write(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// Sendemail 发送验证码邮件-登录
// @Summary 发送验证码
// @Description 接收用户的邮箱地址后判断该用户是否注册并发送验证码邮件
// @Tags 用户登录
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 200 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 404 {object} models.Response "该用户未注册"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/login/sendemail [post]
func (u *UserControllers) LoginSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.Login(email.Email); err != nil {
		c.JSON(http.StatusNotFound, models.Write(err.Error()))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		response := models.Write(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// VerifyEmail 验证验证码-注册
// @Summary 验证验证码
// @Description 接收用户提交的验证码并验证其有效性
// @Tags 用户注册
// @Accept application/json
// @Produce application/json
// @Param code body models.Code true "验证码信息"
// @Success 200 {object} models.Response "验证码验证成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "检验验证码失败"
// @Router /api/register/sendemail/verifycode [post]
func (u *UserControllers) ReigisterVerifyEmail(c *gin.Context) {
	var code models.Code
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		response := err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码验证成功！"))

}

// VerifyEmail 验证验证码-登录
// @Summary 验证验证码
// @Description 接收用户提交的验证码并验证其有效性,并返回登录信息
// @Tags 用户登录
// @Accept application/json
// @Produce application/json
// @Param code body models.Emailverify true "验证码和邮箱信息"
// @Success 200 {object} map[string]interface{} "验证码验证成功并返回token和用户信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "检验验证码失败"
// @Router /api/login/sendemail/verifycode [post]
func (u *UserControllers) LoginVerifyEmail(c *gin.Context) {
	var code models.Emailverify
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		response := err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	User, err := u.EmailService.UserDao.GetUserFromEmail(code.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
	}
	token := u.TokenService.GenerateToken(User)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功!", "token": token, "usermessage": *User})

}

// Register 用户注册
// @Summary 用户注册
// @Description 接收用户的注册信息，生成默认用户名并完成注册，返回用户信息和认证 token
// @Tags 用户注册
// @Accept application/json
// @Produce application/json
// @Param user body models.User true "用户注册信息"
// @Success 201 {object} map[string]interface{} "注册成功的响应，返回token和用户信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "数据库错误或其他服务器错误"
// @Router /api/register/sendemail/verifycode/sendmessage [post]
func (u *UserControllers) Register(c *gin.Context) {
	var User models.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
	}
	//设置身份默认为用户
	User.Role = "user"
	//随机用户uid
	User.ID = u.EmailService.RandomCode(9)
	//防止id重复
	for {
		err := u.UserService.UserDao.SearchUserid(User.ID)
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	//生成随机默认用户名
	rand := u.EmailService.RandomCode(5)
	User.Username = "小知用户" + rand
	if err := u.UserService.Register(&User); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	token := u.TokenService.GenerateToken(&User)
	//成功响应后返回token和用户信息
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "token": token, "userMessage": User})
}

// LoginByPassword 用户登录
// @Summary 密码登录
// @Description 通过密码和邮箱进行登录，检验成功后返回用户信息和认证 token
// @Tags 用户登录
// @Accept application/json
// @Produce application/json
// @Param user body models.Login true "用户注册信息"
// @Success 201 {object} map[string]interface{} "注册成功的响应，返回token和用户信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 404 {object} models.Response "该用户未找到"
// @Failure 409 {object} models.Response "密码不正确"
// @Router /api/login/bypassword [post]
func (u *UserControllers) LoginByPassword(c *gin.Context) {
	var User models.Login
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Usermessage, err := u.UserService.LoginByPassword(User)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Write(err.Error()))
		return
	}
	if err := u.UserService.ComparePassword((*Usermessage).Password, User.Password); err != nil {
		c.JSON(http.StatusConflict, models.Write("密码不正确"))
		return
	}
	token := u.TokenService.GenerateToken(Usermessage)
	c.JSON(http.StatusOK, gin.H{"message": "登陆成功！", "token": token, "usermessage": *Usermessage})
}

// UpdateSendemail 发送验证码邮件-忘记密码
// @Summary 发送验证码
// @Description 接收用户的邮箱地址后发送验证码邮件
// @Tags 忘记密码
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 200 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/forgetPassword/sendemail [post]
func (u *UserControllers) UpdateSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// VerifyEmail 验证验证码-忘记密码
// @Summary 验证验证码
// @Description 接收用户提交的验证码并验证其有效性
// @Tags 忘记密码
// @Accept application/json
// @Produce application/json
// @Param code body models.Code true "验证码信息"
// @Success 200 {object} models.Response "验证码验证成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "检验验证码失败"
// @Router /api/forgetPassword/sendemail/verifycode [post]
func (u *UserControllers) UpdateVerifyEmail(c *gin.Context) {
	var code models.Emailverify
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码验证成功！"))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 通过邮箱进行修改密码
// @Tags 忘记密码
// @Accept application/json
// @Produce application/json
// @Param user body models.Login true "邮箱和新密码"
// @Success 201 {object} models.Response "修改密码成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "修改密码失败"
// @Router /api/forgetPassword/sendemail/verifycode/update [post]
func (u *UserControllers) ChangePassword(c *gin.Context) {
	var User models.Login
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.UserService.ChangePassword(User); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("密码修改成功！"))
}
