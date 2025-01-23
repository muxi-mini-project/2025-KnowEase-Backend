package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeControllers struct {
	LikeService *services.LikeService
	PostService *services.PostService
	UserService *services.UserService
}

func NewLikeControllers(LikeService *services.LikeService, PostService *services.PostService, UserService *services.UserService) *LikeControllers {
	return &LikeControllers{LikeService: LikeService, PostService: PostService, UserService: UserService}
}

// @Summary 用户点赞帖子
// @Description 用户给帖子点赞，点赞数加一
// @Tags 帖子-点赞
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户点赞成功"
// @Failure 500 {object} models.Response "点赞记录上传失败"
// @Router /api/{userid}/post/{postid}/like [post]
func (lc *LikeControllers) LikePost(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	Post, _ := lc.PostService.PostDao.SearchPostByID(PostID)
	User, _ := lc.UserService.UserDao.GetUserFromID(UserID)
	if err := lc.LikeService.LikePost(PostID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	message := fmt.Sprintf("用户%s点赞了你的帖子%s!", User.Username, Post.Title)
	if err := lc.LikeService.InitMessage(Post.PosterID, message, User.ImageURL); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("点赞处理成功，消息上次失败"))
	}
	c.JSON(http.StatusCreated, models.Write("用户点赞成功！"))
}

// @Summary 用户取消点赞帖子
// @Description 用户取消点赞，点赞数减一
// @Tags 帖子-点赞
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户D"
// @Success 201 {object} models.Response "用户取消点赞成功"
// @Failure 500 {object} models.Response "取消点赞记录上传失败"
// @Router /api/{userid}/post/{postid}/like [post]
func (lc *LikeControllers) CancelLike(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if err := lc.LikeService.CancelLike(PostID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消点赞成功！"))
}

// @Summary 获取用户的总获赞数
// @Description 用户给帖子点赞，点赞数加一
// @Tags 个人主页
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} models.Response "查询用户获赞数成功"
// @Failure 500 {object} models.Response "查询用户获赞数失败"
// @Router /api/{userid}/userpage/likecount [get]
func (lc *LikeControllers) GetUserLikes(c *gin.Context) {
	UserID := c.Param("userid")
	count, err := lc.LikeService.GetUserLikes(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询用户获赞数成功！", "UserLikesCount": count})
}

// @Summary 获取用户历史一个月内的历史点赞记录
// @Description 在个人主页获取用户一个月（31天）内的点赞记录
// @Tags 个人主页-记录查询
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功信息以及历史点赞帖子"
// @Failure 500 {object} models.Response "获取点赞记录失败"
// @Router /api/{userid}/userpage/likerecord [get]
func (lc *LikeControllers) GetLikeRecord(c *gin.Context) {
	UserID := c.Param("userid")
	Record, err := lc.LikeService.GetLikeRecord(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	PostMessages, Err := lc.PostService.SearchPostByID(Record)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("点赞记录查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "查询历史点赞记录成功！",
		"LikeRecord": PostMessages,
	})
}

// @Summary 用户收藏帖子
// @Description 用户给帖子收藏，收藏数加一
// @Tags 帖子-收藏
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户收藏成功"
// @Failure 500 {object} models.Response "收藏记录上传失败"
// @Router /api/posts/{postid}/{userid}/save [post]
func (lc *LikeControllers) SavePost(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if err := lc.LikeService.SavePost(PostID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("用户收藏成功！"))
}

// @Summary 用户取消收藏帖子
// @Description 用户取消收藏，收藏数减一
// @Tags 帖子-收藏
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户取消收藏成功"
// @Failure 500 {object} models.Response "取消收藏记录上传失败"
// @Router /api/{userid}/post/{postid}/cancelsave [post]
func (lc *LikeControllers) CancelSave(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if err := lc.LikeService.CancelSave(PostID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消收藏成功！"))
}

// @Summary 获取用户历史收藏记录
// @Description 在个人主页获取用户的历史收藏记录
// @Tags 个人主页-记录查询
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功信息以及历史收藏帖子"
// @Failure 500 {object} models.Response "获取收藏记录失败"
// @Router /api/{userid}/userpage/saverecord [get]
func (lc *LikeControllers) GetSaveRecord(c *gin.Context) {
	UserID := c.Param("userid")
	Record, err := lc.LikeService.GetSaveRecord(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	var Records []string
	for _, Post := range Record {
		Records = append(Records, Post.PostID)
	}
	PostMessages, Err := lc.PostService.SearchPostByID(Records)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("收藏记录查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "查询历史收藏记录成功！",
		"LikeRecord": PostMessages,
	})
}

// @Summary 获取用户浏览记录
// @Description 在个人主页获取用户的历史浏览记录
// @Tags 个人主页-记录查询
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功信息以及历史浏览帖子"
// @Failure 500 {object} models.Response "获取浏览记录失败"
// @Router /api/{userid}/userpage/viewrecord [get]
func (lc *LikeControllers) GetViewRecord(c *gin.Context) {
	UserID := c.Param("userid")
	Record, err := lc.LikeService.GetViewRecord(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	var Records []string
	for _, Post := range Record {
		Records = append(Records, Post.PostID)
	}
	PostMessages, Err := lc.PostService.SearchPostByID(Records)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("用户浏览记录查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "查询历史浏览记录成功！",
		"LikeRecord": PostMessages,
	})
}

// @Summary 用户点赞评论
// @Description 用户给评论点赞，评论数加一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param commentid path string true "评论ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户收藏成功"
// @Failure 500 {object} models.Response "收藏记录上传失败"
// @Router /api/{userid}/post/{postid}/{commentid}/like [post]
func (lc *LikeControllers) LikeComment(c *gin.Context) {
	CommentID := c.Param("commentid")
	UserID := c.Param("userid")
	if err := lc.LikeService.LikeComment(CommentID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("点赞评论成功！"))
}

// @Summary 用户取消点赞评论
// @Description 用户取消点赞评论，评论数减一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param commentid path string true "评论ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户取消点赞成功"
// @Failure 500 {object} models.Response "取消点赞记录上传失败"
// @Router /api/{userid}/post/{postid}/{commentid}/cancellike [post]
func (lc *LikeControllers) CancelCommentLike(c *gin.Context) {
	CommentID := c.Param("commentid")
	UserID := c.Param("userid")
	if err := lc.LikeService.CancelCommentLike(CommentID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消点赞成功！"))
}

// @Summary 用户点赞回复
// @Description 用户给回复点赞，评论数加一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param replyid path string true "回复ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户点赞成功"
// @Failure 500 {object} models.Response "点赞记录上传失败"
// @Router /api/{userid}/post/{postid}/{commentid}/{replyid}/like [post]
func (lc *LikeControllers) LikeReply(c *gin.Context) {
	ReplyID := c.Param("replyid")
	UserID := c.Param("userid")
	if err := lc.LikeService.LikeReply(ReplyID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("点赞评论成功！"))
}

// @Summary 用户取消点赞回复
// @Description 用户取消点赞回复，点赞数减一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param replyid path string true "回复ID"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "用户取消点赞成功"
// @Failure 500 {object} models.Response "取消点赞记录上传失败"
// @Router /api/{userid}/post/{postid}/{commentid}/{replyid}/cancellike [post]
func (lc *LikeControllers) CancelReplyLike(c *gin.Context) {
	ReplyID := c.Param("replyid")
	UserID := c.Param("userid")
	if err := lc.LikeService.CancelReplyLike(ReplyID, UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消点赞成功！"))
}

// 定期更新数据至数据库
func (lc *LikeControllers) UpdateAllCount() {
	lc.LikeService.StartUpdateTicker()
}
