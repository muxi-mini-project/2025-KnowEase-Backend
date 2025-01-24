package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type UserSvc struct {
	uc *controllers.UserControllers
	m  *middleware.Middleware
}

func NewUserSvc(uc *controllers.UserControllers, m *middleware.Middleware) *UserSvc {
	return &UserSvc{
		uc: uc,
		m:  m,
	}
}
func (u *UserSvc) NewUserGroup(r *gin.Engine) {
	r.Use(u.m.Cors())
	registerGroup := r.Group("/api/register")
	{
		registerGroup.POST("/sendemail", u.uc.RegisterSendemail)
		registerGroup.POST("/sendemail/verifycode", u.uc.ReigisterVerifyEmail)
		registerGroup.POST("/sendemail/verifycode/sendmessage", u.uc.Register)
	}
	loginGroup := r.Group("/api/login")
	{
		loginGroup.POST("/bypassword", u.uc.LoginByPassword)
		loginGroup.POST("/sendemail", u.uc.LoginSendemail)
		loginGroup.POST("/sendemail/verifycode", u.uc.LoginVerifyEmail)
	}
	forgetPassword := r.Group("/api/forgetPassword")
	{
		forgetPassword.POST("/sendemail", u.uc.UpdateSendemail)
		forgetPassword.POST("/sendemail/verifycode", u.uc.UpdateVerifyEmail)
		forgetPassword.POST("/sendemail/verifycode/update", u.uc.ChangePassword)
	}

}
