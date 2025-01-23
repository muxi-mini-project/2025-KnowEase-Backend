package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type UserPageSvc struct {
	lc *controllers.LikeControllers
	m  *middleware.Middleware
	pc *controllers.PostControllers
}

func NewUserPageSvc(lc *controllers.LikeControllers, m *middleware.Middleware, pc *controllers.PostControllers) *UserPageSvc {
	return &UserPageSvc{
		lc: lc,
		m:  m,
		pc: pc,
	}
}
func (up *UserPageSvc) NewUserPageGroup(r *gin.Engine) {
	r.Use(up.m.Cors())
	r.Use(up.m.Verifytoken())
	userpage := r.Group("/api")
	{
		userpage.GET("/:userid/userpage/likecount", up.lc.GetUserLikes)
		userpage.GET("/:userid/userpage/likerecord", up.lc.GetLikeRecord)
		userpage.GET("/:userid/userpage/viewrecord", up.lc.GetViewRecord)
		userpage.GET("/:userid/userpage/saverecord", up.lc.GetSaveRecord)
		userpage.DELETE("/:userid/userpage/mypost/delete/:postid", up.pc.DeletePost)
		userpage.DELETE("/:userid/userpage/mypost/deleteposts", up.pc.DeletePosts)
		userpage.GET("/:userid/userpage/message", up.pc.GetUserUnreadMessage)
	}
}
