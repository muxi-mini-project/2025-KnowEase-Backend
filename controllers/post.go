package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"fmt"
	"log"

	"net/http"

	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostControllers struct {
	PostService  *services.PostService
	LikeService  *services.LikeService
	EmailService *services.EmailService
	UserService  *services.UserService
}

func NewPostControllers(PostService *services.PostService, LikeService *services.LikeService, EmailService *services.EmailService, UserService *services.UserService) *PostControllers {
	return &PostControllers{PostService: PostService, LikeService: LikeService, EmailService: EmailService, UserService: UserService}
}

// @Summary 用户发帖
// @Description 用户发帖，上传帖子正文、图片等信息。
// @Tags 生活区-发帖
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Param post body models.PostMessage true "发帖内容"
// @Success 201 {object} map[string]interface{} "成功的响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询发帖人信息失败！"
// @Failure 500 {object} models.Response "帖子发布失败！"
// @Router /api/{userid}/post/publish [post]
func (pc *PostControllers) PublishPostBody(c *gin.Context) {
	var Post models.PostMessage
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	Post.PosterID = UserID
	var URLs []string
	if err := json.Unmarshal([]byte(Post.ImageURL), &URLs); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("URL输入无效，请重试！"))
		return
	}
	//拼接url并用逗号分隔开
	imageURL := strings.Join(URLs, ",")
	Post.ImageURL = imageURL
	for {
		//生成帖子id
		Post.PostID = pc.EmailService.RandomCode(6)
		err := pc.PostService.PostDao.SearchPostID(Post.PostID)
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := pc.PostService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询发帖人信息失败！"))
		return
	}
	Post.PosterName = PosterName
	Post.PosterURL = PosterURL
	if err := pc.PostService.PublishPost(Post); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("帖子发布失败！"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发帖成功！", "PostMessage": Post})
}

// @Summary 删除帖子
// @Description 用户在个人主页选择删除帖子
// @Tags 个人主页-我的发布
// @Accept  json
// @Produce  json
// @Param postid path string true "帖子ID"
// @Success 201 {object} models.Response "删帖成功！"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "删除帖子及其相关信息失败！"
// @Router /api/{userid}/userpage/mypost/delete/{postid} [delete]
func (pc *PostControllers) DeletePost(c *gin.Context) {
	PostID := c.Param("postid")
	if PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	var PostIDs []string
	PostIDs = append(PostIDs, PostID)
	if err := pc.PostService.DeletePost(PostIDs); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除帖子失败！"))
		return
	}
	if err := pc.PostService.DeletePostComment(PostIDs); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除帖子相关评论失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删帖成功！"))
}

// @Summary 批量删除帖子
// @Description 用户在个人主页选择批量删除帖子
// @Tags 个人主页-删帖
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "删帖成功！"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "删除帖子及其相关信息失败！"
// @Router /api/{userid}/userpage/mypost/deleteposts [delete]
func (pc *PostControllers) DeletePosts(c *gin.Context) {
	var Post models.PostIDs
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效,请重试！"))
		return
	}
	var PostIDS []string
	if err := json.Unmarshal([]byte(Post.PostID), &PostIDS); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("ID输入无效，请重试！"))
		return
	}
	if err := pc.PostService.DeletePost(PostIDS); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("批量删除帖子失败！"))
		return
	}
	if err := pc.PostService.DeletePostComment(PostIDS); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除相关帖子评论失败！"))
		return
	}
	if err := pc.PostService.DeletePostComment(PostIDS); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除相关帖子评论失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("批量删除帖子成功！"))
}

// @Summary 推荐标签页面
// @Description 获取生活区的推荐标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "推荐帖子查询失败"
// @Router /api/{userid}/post/recommend [get]
func (pc *PostControllers) RecommendationPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	WeightedRecord, err := pc.PostService.WeightedRecommendation(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	var PostIDs []string
	for _, Posts := range WeightedRecord {
		PostIDs = append(PostIDs, Posts.PostID)
	}
	PostMessages, Err := pc.PostService.SearchPostByID(PostIDs)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("推荐帖子信息查询失败！"))
		return
	}
	for i := 0; i < len(PostMessages); i++ {
		parts := strings.Split(PostMessages[i].ImageURL, ",")
		jsonData, err := json.Marshal(parts)
		if err != nil {
			log.Println("something wrong in imageurls:%w", err)
			continue
		}
		PostMessages[i].ImageURL = string(jsonData)
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求推荐帖成功！", "posts": PostMessages})
}

// @Summary 校园标签页面
// @Description 获取生活区的校园标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "校园帖子查询失败"
// @Router /api/{userid}/post/campus [get]
func (pc *PostControllers) CampusPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "校园")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	for _, Post := range Posts {
		parts := strings.Split(Post.ImageURL, ",")
		jsonData, err := json.Marshal(parts)
		if err != nil {
			log.Println("something wrong in imageurls:%w", err)
			continue
		}
		Post.ImageURL = string(jsonData)
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求校园标签帖成功！", "posts": Posts})
}

// @Summary 生活标签页面
// @Description 获取生活区的生活标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "生活帖子查询失败"
// @Router /api/{userid}/post/life [get]
func (pc *PostControllers) LifePost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "生活")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	for i := 0; i < len(Posts); i++ {
		parts := strings.Split(Posts[i].ImageURL, ",")
		jsonData, err := json.Marshal(parts)
		if err != nil {
			log.Println("something wrong in imageurls:%w", err)
			continue
		}
		Posts[i].ImageURL = string(jsonData)
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求生活标签帖成功！", "posts": Posts})
}

// @Summary 美食标签页面
// @Description 获取生活区的美食标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "美食帖子查询失败"
// @Router /api/{userid}/post/food [get]
func (pc *PostControllers) FoodPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "美食")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	for i := 0; i < len(Posts); i++ {
		parts := strings.Split(Posts[i].ImageURL, ",")
		jsonData, err := json.Marshal(parts)
		if err != nil {
			log.Println("something wrong in imageurls:%w", err)
			continue
		}
		Posts[i].ImageURL = string(jsonData)
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求美食标签帖成功！", "posts": Posts})
}

// @Summary 绘画标签页面
// @Description 获取生活区的绘画标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "绘画帖子查询失败"
// @Router /api/{userid}/post/paint [get]
func (pc *PostControllers) PaintPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "绘画")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	for i := 0; i < len(Posts); i++ {
		parts := strings.Split(Posts[i].ImageURL, ",")
		jsonData, err := json.Marshal(parts)
		if err != nil {
			log.Println("something wrong in imageurls:%w", err)
			continue
		}
		Posts[i].ImageURL = string(jsonData)
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求绘画标签帖成功！", "posts": Posts})
}

// @Summary 发布评论
// @Description 用户在帖子下面发布评论
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Param postid path string true "帖子ID"
// @Success 201 {object} map[string]interface{} "成功响应信息以及评论信息"
// @Failure 500 {object} models.Response "发布失败"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/publishcomment [post]
func (pc *PostControllers) PublishComment(c *gin.Context) {
	var Post models.Comment
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if UserID == "" || PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	Post.PostID = PostID
	Post.CommenterID = UserID
	for {
		//生成评论id
		Post.CommentID = pc.EmailService.RandomCode(6)
		err := pc.PostService.PostDao.SearchCommentID(Post.CommentID)
		err1 := pc.PostService.PostDao.SearchReplyID(Post.CommentID)
		if err == gorm.ErrRecordNotFound && err1 == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := pc.PostService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询回复人信息失败！"))
		return
	}
	Post.CommenterName = PosterName
	Post.CommenterURL = PosterURL
	if err := pc.PostService.PublishComment(Post); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("评论发布失败！"))
		return
	}
	PostMessage, _ := pc.PostService.PostDao.SearchPostByID(PostID)
	User, _ := pc.UserService.UserDao.GetUserFromID(UserID)
	message := fmt.Sprintf("用户%s评论了帖子%s!", User.Username, PostMessage.Title)
	if err := pc.LikeService.InitMessage(PostMessage.PosterID, message, User.ImageURL); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("回复处理成功，消息上次失败"))
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发布评论成功！", "CommentMessage": Post})
}

// @Summary 删除评论
// @Description 用户在帖子下面发布评论
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param commentid path string true "评论ID"
// @Success 201 {object} models.Response "删除成功"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/post/{postid}/deletecomment [delete]
func (pc *PostControllers) DeleteComment(c *gin.Context) {
	CommentID := c.Param("commentid")
	if CommentID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := pc.PostService.DeleteComment(CommentID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论失败！"))
		return
	}
	if err := pc.PostService.DeleteAllReply(CommentID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论回复失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删除评论成功！"))
}

// @Summary 发布回复
// @Description 用户在帖子下面发布回复
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Param postid path string true "帖子ID"
// @Success 207 {object} models.Response "发布成功但消息发送失败"
// @Success 201 {object} map[string]interface{} "成功响应信息以及评论信息"
// @Failure 500 {object} models.Response "查询失败"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/{commentid}/publishreply [post]
func (pc *PostControllers) PublishReply(c *gin.Context) {
	var Post models.Reply
	CommentID, UserID, PostID := c.Param("comment_id"), c.Param("userid"), c.Param("postid")
	if CommentID == "" || UserID == "" || PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	Post.CommentID = CommentID
	Post.ReplyerID = UserID
	for {
		//生成回复id
		Post.ReplyID = pc.EmailService.RandomCode(6)
		err := pc.PostService.PostDao.SearchCommentID(Post.ReplyID)
		err1 := pc.PostService.PostDao.SearchReplyID(Post.ReplyID)
		if err == gorm.ErrRecordNotFound && err1 == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := pc.PostService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询回复人信息失败！"))
		return
	}
	Post.ReplyerName = PosterName
	Post.ReplyURL = PosterURL
	if err := pc.PostService.PublishReply(Post); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("回复发布失败！"))
		return
	}
	PostMessage, _ := pc.PostService.PostDao.SearchPostByID(PostID)
	CommentMessage, _ := pc.PostService.PostDao.SearchCommentByID(CommentID)
	User, _ := pc.UserService.UserDao.GetUserFromID(UserID)
	message := fmt.Sprintf("用户%s回复了你在帖子%s的评论!", User.Username, PostMessage.Title)
	if err := pc.LikeService.InitMessage(CommentMessage.CommenterID, message, User.ImageURL); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("回复处理成功，消息上次失败"))
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发布回复成功！", "ReplyMessage": Post})
}

// @Summary 删除回复
// @Description 用户删除自己的回复
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param replyid path string true "回复ID"
// @Success 201 {object} models.Response "删除成功"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "删除失败"
// @Router /api/{userid}/post/{postid}/{commentid}/{replyid} [delete]
func (pc *PostControllers) DeleteReply(c *gin.Context) {
	ReplyID := c.Param("replyid")
	if ReplyID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := pc.PostService.DeleteReply(ReplyID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论失败！"))
		return
	}
	//删除回复的回复
	if err := pc.PostService.DeleteAllReply(ReplyID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论回复失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删除评论成功！"))
}

// @Summary 获取帖子详细信息
// @Description 用户点开帖子详情页，获取帖子详情信息,帖子浏览量加一
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param postid path string true "帖子ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/post/{postid} [get]
func (pc *PostControllers) GetPostMessage(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if PostID == "" || UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	pc.LikeService.ViewPost(PostID, UserID)
	PostMessage, err := pc.PostService.GetAllComment(PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询帖子具体信息出错！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询帖子信息成功！", "postMessage": *PostMessage})
}

// @Summary 获取未读消息通知
// @Description 用户点开消息通知，获取未读的消息（被点赞，被评论），并把消息状态更新为已读
// @Tags 个人主页-消息通知
// @Accept  json
// @Produce  json
// @Param userid path string true "帖子ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及消息信息"
// @Success 207 {object} map[string]interface{} "状态更新错误以及消息信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/userpage/message [get]
func (pc *PostControllers) GetUserUnreadMessage(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Message, err := pc.PostService.SearchAllUnreadMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询未读消息出错！"))
		return
	}
	//将所有消息更新为已读
	err = pc.PostService.UpdateMessageStatus(UserID)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{
			"message":  "消息更新出错",
			"messages": Message,
		})
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取并更新消息成功！", "messages": Message})
}
