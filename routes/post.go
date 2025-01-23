package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type PostSvc struct {
	lc *controllers.LikeControllers
	pc *controllers.PostControllers
	m  *middleware.Middleware
}

func NewPostSvc(lc *controllers.LikeControllers, m *middleware.Middleware, pc *controllers.PostControllers) *PostSvc {
	return &PostSvc{
		lc: lc,
		m:  m,
		pc: pc,
	}
}
func (p *PostSvc) NewPostGroup(r *gin.Engine) {
	go p.lc.UpdateAllCount()
	r.Use(p.m.Cors())
	r.Use(p.m.Verifytoken())
	posts := r.Group("/api")
	{
		posts.POST("/:userid/post/publish", p.pc.PublishPostBody)
		posts.GET("/:userid/post/recommend", p.pc.RecommendationPost)
		posts.GET("/:userid/post/campus", p.pc.CampusPost)
		posts.GET("/:userid/post/life", p.pc.LifePost)
		posts.GET("/:userid/post/food", p.pc.FoodPost)
		posts.GET("/:userid/post/paint", p.pc.PaintPost)
		posts.GET("/:userid/post/:postid", p.pc.GetPostMessage)
		posts.POST("/:userid/post/:postid/publishcomment", p.pc.PublishComment)
		posts.DELETE("/:userid/post/:postid/deletecomment", p.pc.DeleteComment)
		posts.POST("/:userid/post/:postid/:commentid/publishreply", p.pc.PublishReply)
		posts.DELETE("/:userid/post/:postid/:commentid/:replyid", p.pc.DeleteReply)
		posts.POST("/:userid/post/:postid/like", p.lc.LikePost)
		posts.POST("/:userid/post/:postid/cancellike", p.lc.CancelLike)
		posts.POST("/:userid/post/:postid/save", p.lc.SavePost)
		posts.POST("/:userid/post/:postid/cancelsave", p.lc.CancelSave)
		posts.POST("/:userid/post/:postid/:commentid/like", p.lc.LikeComment)
		posts.POST("/:userid/post/:postid/:commentid/cancellike", p.lc.CancelCommentLike)
		posts.POST("/:userid/post/:postid/:commentid/:replyid/like", p.lc.LikeReply)
		posts.POST("/:userid/post/:postid/:commentid/:replyid/cancellike", p.lc.CancelReplyLike)
		posts.GET("/:userid/post/:postid/getcounts", p.lc.GetPostCounts)
		posts.GET("/:userid/post/:postid/:commentid/getcounts", p.lc.GetCommentCounts)
		posts.GET("/:userid/post/:postid/:commentid/:replyid/getcounts", p.lc.GetReplyCounts)
	}
}
