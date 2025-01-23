package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"KnowEase/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type LikeService struct {
	ld *dao.LikeDao
	pd *dao.PostDao
}

func NewLikeService(ld *dao.LikeDao, pd *dao.PostDao) *LikeService {
	return &LikeService{ld: ld, pd: pd}
}

// 用户点赞帖子
func (ls *LikeService) LikePost(PostID, UserID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Post:%s:like_users", PostID)
	//查找点赞记录以避免重复点赞
	if utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is already liked post %s", UserID, PostID)
	}
	//添加用户到点赞集合
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新用户点赞记录
	Record := models.UserLikeHistory{
		UserID: UserID,
		PostID: PostID,
	}
	go ls.ld.SyncUserLikeHistoryToDB(&Record)
	//更新帖子总获赞数
	key = fmt.Sprintf("Poster:%s:likes", PosterID.PosterID)
	utils.Client.HIncrBy(utils.Ctx, key, PostID, 1)
	return nil
}

// 用户取消点赞
func (ls *LikeService) CancelLike(PostID, UserID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Post:%s:like_users", PostID)
	//查找点赞记录判断点赞状态
	if !utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s has not liked post %s", UserID, PostID)
	}
	//将用户从点赞集合中移除
	utils.Client.SRem(utils.Ctx, key, UserID)
	go ls.ld.DeleteUserLikeHistory(PostID)
	//更新帖子总获赞数
	key = fmt.Sprintf("Poster:%s:likes", PosterID.PosterID)
	utils.Client.HIncrBy(utils.Ctx, key, PostID, -1)
	return nil
}

// 获取帖子的总点赞数
func (ls *LikeService) GetPostLikes(PostID string) (int, error) {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return 0, fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Poster:%s:likes", PosterID.PosterID)
	likescount, err := utils.Client.HGet(utils.Ctx, key, PostID).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(likescount)
	go ls.ld.SyncPostLikeToDB(PostID, count)
	return count, err
}

// 获取用户总获赞数
func (ls *LikeService) GetUserLikes(UserID string) (int, error) {
	key := fmt.Sprintf("Post:%s:likes", UserID)
	likescount, err := utils.Client.HGetAll(utils.Ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	// 遍历累加点赞数
	var count, totalLikes int
	count = 0
	totalLikes = 0
	for _, countStr := range likescount {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			log.Println("failed to count this record:%w", err)
		}
		totalLikes += count
	}
	//更新用户总获赞数
	go ls.ld.SyncUserLikesToDB(UserID, totalLikes)
	return totalLikes, nil
}

// 查询用户历史点赞记录
func (ls *LikeService) GetLikeRecord(UserID string) ([]string, error) {
	UserHistory, err := ls.ld.GetLikeRecord(UserID)
	if err != nil {
		return nil, err
	}
	var Record []string
	for _, Message := range UserHistory {
		Record = append(Record, Message.PostID)
	}
	return Record, nil
}

// 用户收藏帖子
func (ls *LikeService) SavePost(PostID, UserID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Post:%s:save_users", PostID)
	//查找收藏记录以避免重复收藏
	if utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is already saved post %s", UserID, PostID)
	}
	//添加用户到收藏集合
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新用户收藏记录
	Record := models.UserSaveHistory{
		UserID: UserID,
		PostID: PostID,
	}
	go ls.ld.SyncUserSaveHistoryToDB(&Record)
	//更新帖子总收藏数
	key = fmt.Sprintf("Poster:%s:saves", PosterID.PosterID)
	utils.Client.HIncrBy(utils.Ctx, key, PostID, 1)
	return nil
}

// 用户取消收藏
func (ls *LikeService) CancelSave(PostID, UserID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Post:%s:save_users", PostID)
	//查找点赞记录判断收藏状态
	if !utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s has not saved post %s", UserID, PostID)
	}
	//将用户从收藏集合中移除
	utils.Client.SRem(utils.Ctx, key, UserID)
	go ls.ld.DeleteUserSaveHistory(PostID)
	//更新帖子总收藏数
	key = fmt.Sprintf("Poster:%s:saves", PosterID.PosterID)
	utils.Client.HIncrBy(utils.Ctx, key, PostID, -1)
	return nil
}

// 获取帖子的总收藏数
func (ls *LikeService) GetPostSaves(PostID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Poster:%s:saves", PosterID.PosterID)
	savescount, err := utils.Client.HGet(utils.Ctx, key, PostID).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	count, err := strconv.Atoi(savescount)
	go ls.ld.SyncPostSaveToDB(PostID, count)
	return err
}

// 查询用户收藏记录
func (ls *LikeService) GetSaveRecord(UserID string) ([]models.UserSaveHistory, error) {
	UserHistory, err := ls.ld.GetSaveRecord(UserID)
	if err != nil {
		return nil, err
	}
	return UserHistory, nil
}

// 用户浏览帖子
func (ls *LikeService) ViewPost(PostID, UserID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Post:%s:view_users", PostID)
	//添加用户到浏览集合
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新用户浏览记录
	Record := models.UserViewHistory{
		UserID:   UserID,
		PostID:   PostID,
		CreateAt: time.Now(),
	}
	go ls.ld.SyncUserViewHistoryToDB(&Record)
	//更新帖子总浏览数
	key = fmt.Sprintf("Poster:%s:views", PosterID.PosterID)
	utils.Client.HIncrBy(utils.Ctx, key, PostID, 1)
	return nil
}

// 获取帖子的总浏览数
func (ls *LikeService) GetPostViews(PostID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Poster:%s:views", PosterID.PosterID)
	viewscount, err := utils.Client.HGet(utils.Ctx, key, PostID).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	count, err := strconv.Atoi(viewscount)
	go ls.ld.SyncPostViewToDB(PostID, count)
	return err
}

// 查询用户浏览记录
func (ls *LikeService) GetViewRecord(UserID string) ([]models.UserViewHistory, error) {
	UserHistory, err := ls.ld.GetViewRecord(UserID)
	if err != nil {
		return nil, err
	}
	return UserHistory, nil
}

// 创建消息推送
func (ls *LikeService) InitMessage(UserID, message, Userurl string) error {
	var Message models.Message
	Message.Status = "unread"
	Message.UserID = UserID
	Message.Message = message
	Message.PosterURL = Userurl
	return ls.ld.SyncMessageToDB(&Message)
}

// 用户点赞评论
func (ls *LikeService) LikeComment(CommentID, UserID string) error {
	key := fmt.Sprintf("comment:%s:like_users", CommentID)
	//查找点赞记录以避免重复点赞
	if utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is already liked comment %s", UserID, CommentID)
	}
	//添加用户到点赞集合
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新评论总点赞数
	key = "Comment:likes"
	utils.Client.HIncrBy(utils.Ctx, key, CommentID, 1)
	return nil
}

// 用户取消评论点赞
func (ls *LikeService) CancelCommentLike(CommentID, UserID string) error {
	key := fmt.Sprintf("comment:%s:like_users", CommentID)
	//查找点赞记录判断点赞状态
	if !utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s has not liked post %s", UserID, CommentID)
	}
	//将用户从点赞集合中移除
	utils.Client.SRem(utils.Ctx, key, UserID)
	//更新帖子总收藏数
	key = "Comment:likes"
	utils.Client.HIncrBy(utils.Ctx, key, CommentID, -1)
	return nil
}

// 获取评论的总点赞数
func (ls *LikeService) GetCommentLikes(CommentID string) (int, error) {
	key := "Comment:likes"
	savescount, err := utils.Client.HGet(utils.Ctx, key, CommentID).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(savescount)
	go ls.ld.SyncCommentLikeToDB(CommentID, count)
	return count, err
}

// 用户点赞回复
func (ls *LikeService) LikeReply(ReplyID, UserID string) error {
	key := fmt.Sprintf("reply:%s:like_users", ReplyID)
	//查找点赞记录以避免重复点赞
	if utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is already liked reply %s", UserID, ReplyID)
	}
	//添加用户到点赞集合
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新评论总点赞数
	key = "Reply:likes"
	utils.Client.HIncrBy(utils.Ctx, key, ReplyID, 1)
	return nil
}

// 用户取消评论点赞
func (ls *LikeService) CancelReplyLike(ReplyID, UserID string) error {
	key := fmt.Sprintf("reply:%s:like_users", ReplyID)
	//查找点赞记录判断收藏状态
	if !utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s has not liked reply %s", UserID, ReplyID)
	}
	//将用户从点赞集合中移除
	utils.Client.SRem(utils.Ctx, key, UserID)
	//更新帖子总收藏数
	key = "Reply:likes"
	utils.Client.HIncrBy(utils.Ctx, key, ReplyID, -1)
	return nil
}

// 获取评论的总点赞数
func (ls *LikeService) GetReplyLikes(ReplyID string) (int, error) {
	key := "Reply:likes"
	savescount, err := utils.Client.HGet(utils.Ctx, key, ReplyID).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(savescount)
	go ls.ld.SyncReplyLikeToDB(ReplyID, count)
	return count, err
}

// 设置定时更新数据库数据
func (ls *LikeService) StartUpdateTicker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ls.UpdateAllCount()
		}
	}
}

func (ls *LikeService) UpdateAllCount() {
	Posts, _ := ls.pd.SearchAllPost()
	for _, Post := range Posts {
		ls.GetPostViews(Post.PostID)
		ls.GetPostSaves(Post.PostID)
		ls.GetPostSaves(Post.PostID)
	}
	Comments, _ := ls.pd.SearchAllComments()
	for _, Comment := range Comments {
		ls.GetCommentLikes(Comment.CommentID)
	}
	Replys, _ := ls.pd.SearchALLReplys()
	for _, Reply := range Replys {
		ls.GetReplyLikes(Reply.ReplyID)
	}
}
