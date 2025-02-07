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
	ld dao.LikeDaoInterface
	pd dao.PostDaoInterface
	ud dao.UserDaoInterface
}

func NewLikeService(ld dao.LikeDaoInterface, pd dao.PostDaoInterface, ud dao.UserDaoInterface) *LikeService {
	return &LikeService{ld: ld, pd: pd, ud: ud}
}

// 用户点赞帖子
func (ls *LikeService) LikePost(PostID, UserID string) error {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Post:%s:like_users", PostID)
	fmt.Print(key)
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
	return count, err
}

// 获取用户总获赞数
func (ls *LikeService) GetUserLikes(UserID string) (int, error) {
	key := fmt.Sprintf("Poster:%s:likes", UserID)
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
func (ls *LikeService) GetPostSaves(PostID string) (int, error) {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return 0, fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Poster:%s:saves", PosterID.PosterID)
	savescount, err := utils.Client.HGet(utils.Ctx, key, PostID).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(savescount)
	return count, err
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
func (ls *LikeService) GetPostViews(PostID string) (int, error) {
	PosterID, err := ls.ld.SearchPostmessageFromPost(PostID)
	if err != nil {
		return 0, fmt.Errorf("failed to find the poster:%w", err)
	}
	key := fmt.Sprintf("Poster:%s:views", PosterID.PosterID)
	viewscount, err := utils.Client.HGet(utils.Ctx, key, PostID).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(viewscount)
	return count, err
}

// 将帖子浏览数写入数据库
func (ls *LikeService) SyncPostViewToDB(PostID string) {
	count, err := ls.GetPostViews(PostID)
	if err != nil {
		log.Printf("update post %s viewcount error:%v", PostID, err)
	}
	err = ls.ld.SyncPostViewToDB(PostID, count)
	if err != nil {
		log.Printf("update post %s viewcount error:%v", PostID, err)
	}
}

// 将帖子点赞数写入数据库
func (ls *LikeService) SyncPostLikeToDB(PostID string) {
	count, err := ls.GetPostLikes(PostID)
	if err != nil {
		log.Printf("update post %s likecount error:%v", PostID, err)
	}
	err = ls.ld.SyncPostLikeToDB(PostID, count)
	if err != nil {
		log.Printf("update post %s likecount error:%v", PostID, err)
	}
}

// 将帖子收藏数写入数据库
func (ls *LikeService) SyncPostSaveToDB(PostID string) {
	count, err := ls.GetPostSaves(PostID)
	if err != nil {
		log.Printf("update post %s savecount error:%v", PostID, err)
	}
	err = ls.ld.SyncPostSaveToDB(PostID, count)
	if err != nil {
		log.Printf("update post %s savecount error:%v", PostID, err)
	}
}

// 将评论点赞数写入数据库
func (ls *LikeService) SyncCommentLikeToDB(CommentID string) {
	count, err := ls.GetCommentLikes(CommentID)
	if err != nil {
		log.Printf("update comment %s likecount error:%v", CommentID, err)
	}
	err = ls.ld.SyncCommentLikeToDB(CommentID, count)
	if err != nil {
		log.Printf("update comment %s likecount error:%v", CommentID, err)
	}
}

// 将帖子评论数写入数据库
func (ls *LikeService) SyncReplyLikeToDB(ReplyID string) {
	count, err := ls.GetReplyLikes(ReplyID)
	if err != nil {
		log.Printf("update reply %s likecount error:%v", ReplyID, err)
	}
	err = ls.ld.SyncReplyLikeToDB(ReplyID, count)
	if err != nil {
		log.Printf("update reply %s likecount error:%v", ReplyID, err)
	}
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
func (ls *LikeService) InitMessage(UserID, message, Userurl, Tag, PostID string) error {
	var Message models.Message
	Message.Status = "unread"
	Message.UserID = UserID
	Message.Message = message
	Message.PosterURL = Userurl
	Message.PostID = PostID
	Message.Tag = Tag
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

// 获取回复的总点赞数
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
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ls.UpdateAllCount()
		}
	}
}

// 获取用户粉丝列表
func (ls *LikeService) GetUserFolower(UserID string) ([]string, error) {
	return ls.ld.SearchUserFollowee(UserID)
}

// 用户关注操作
func (ls *LikeService) Follow(UserID, FolloweeID string) error {
	key := fmt.Sprintf("User:%s:followers", FolloweeID)
	//查找关注记录以避免重复关注
	if utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is already followed user %s", UserID, FolloweeID)
	}
	//添加用户到粉丝合集
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新用户关注数
	key = "User:followees"
	utils.Client.HIncrBy(utils.Ctx, key, UserID, 1)
	Record := models.FollowMessage{
		FollowerID: UserID,
		FolloweeID: FolloweeID,
	}
	go ls.ld.SyncFollowMessageToDB(&Record)
	key = "User:followers"
	utils.Client.HIncrBy(utils.Ctx, key, FolloweeID, 1)
	return nil
}

// 用户取消关注操作
func (ls *LikeService) CancelFollow(UserID, FolloweeID string) error {
	key := fmt.Sprintf("User:%s:followers", FolloweeID)
	//查找关注记录以避免重复关注
	if !utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is not followed user %s", UserID, FolloweeID)
	}
	//用户移除粉丝合集
	utils.Client.SRem(utils.Ctx, key, UserID)
	//更新用户关注数
	key = "User:followees"
	utils.Client.HIncrBy(utils.Ctx, key, UserID, -1)
	go ls.ld.DeleteFollowMessage(UserID)
	key = "User:followers"
	utils.Client.HIncrBy(utils.Ctx, key, FolloweeID, -1)
	return nil
}

// 获取用户关注列表
func (ls *LikeService) GetUserFoloweeList(UserID string) ([]models.Usermessage, error) {
	var FolloweeMessages []models.Usermessage
	UserIDs, err := ls.ld.SearchUserFollowee(UserID)
	if err != nil {
		return nil, err
	}
	for _, FolloweeID := range UserIDs {
		FolloweeMessage, err := ls.ud.GetUserFromID(FolloweeID)
		if err != nil {
			fmt.Printf("get followee%s error:%v", FolloweeID, err)
		} else {
			FolloweeMessages = append(FolloweeMessages, *FolloweeMessage)
		}
	}
	return FolloweeMessages, nil
}

// 获取用户粉丝列表
func (ls *LikeService) GetUserFolowerList(UserID string) ([]models.Usermessage, error) {
	var FolloweeMessages []models.Usermessage
	UserIDs, err := ls.ld.SearchUserFollower(UserID)
	if err != nil {
		return nil, err
	}
	for _, FolloweeID := range UserIDs {
		FolloweeMessage, err := ls.ud.GetUserFromID(FolloweeID)
		if err != nil {
			fmt.Printf("get followee%s error:%v", FolloweeID, err)
		} else {
			FolloweeMessages = append(FolloweeMessages, *FolloweeMessage)
		}
	}
	return FolloweeMessages, nil
}

// 获取用户关注数和粉丝数
func (ls *LikeService) GetFollowCount(UserID string) (int, int, error) {
	var count1, count2 int
	key1 := "User:followers"
	key2 := "User:followees"
	followerscount, err1 := utils.Client.HGet(utils.Ctx, key1, UserID).Result()
	followeescount, err2 := utils.Client.HGet(utils.Ctx, key2, UserID).Result()
	if err1 == redis.Nil {
		count1 = 0
	} else if err2 == redis.Nil {
		count2 = 0
	} else if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("get follow counts error:%v\t%v", err1, err2)
	}
	count1, _ = strconv.Atoi(followerscount)
	count2, _ = strconv.Atoi(followeescount)
	return count1, count2, nil
}

// 关注数据同步
func (ls *LikeService) SyncFollowCount(UserID string) {
	followercount, followeecount, err := ls.GetFollowCount(UserID)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = ls.ld.SyncUserFollowersToDB(UserID, followercount)
	if err != nil {
		fmt.Print("update follower count error%w", err)
	}
	err = ls.ld.SyncUserFolloweesToDB(UserID, followeecount)
	if err != nil {
		fmt.Print("update followee count error%w", err)
	}
}

// 获赞数据同步
func (ls *LikeService) SyncLikeCount(UserID string) {
	totalLikes, err := ls.GetUserLikes(UserID)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = ls.ld.SyncUserLikesToDB(UserID, totalLikes)
	if err != nil {
		fmt.Print("update user like count error%w", err)
	}
}

// 根据id找评论
func (ls *LikeService) SearchCommentByID(CommentID string) (string, string, error) {
	return ls.ld.SearchCommentByID(CommentID)
}

// 根据id找评论
func (ls *LikeService) SearchReplyByID(ReplyID string) (string, string, error) {
	return ls.ld.SearchReplyByID(ReplyID)
}

func (ls *LikeService) UpdateAllCount() {
	Posts, _ := ls.pd.SearchAllPost()
	for _, Post := range Posts {
		ls.SyncPostLikeToDB(Post.PostID)
		ls.SyncPostSaveToDB(Post.PostID)
		ls.SyncPostViewToDB(Post.PostID)
	}
	Comments, _ := ls.pd.SearchAllComments()
	for _, Comment := range Comments {
		ls.SyncCommentLikeToDB(Comment.CommentID)
	}
	Replys, _ := ls.pd.SearchALLReplys()
	for _, Reply := range Replys {
		ls.SyncReplyLikeToDB(Reply.ReplyID)
	}
	Users, _ := ls.ud.SearchAllUser()
	for _, User := range Users {
		ls.SyncFollowCount(User)
		ls.SyncLikeCount(User)
	}
}

// 查询用户点赞帖子状态
func (ls *LikeService) GetPostLikeStatus(UserID, PostID string) bool {
	key := fmt.Sprintf("Post:%s:like_users", PostID)
	return utils.Client.SIsMember(utils.Ctx, key, UserID).Val()
}

// 查询用户收藏帖子状态
func (ls *LikeService) GetPostSaveStatus(UserID, PostID string) bool {
	key := fmt.Sprintf("Post:%s:save_users", PostID)
	return utils.Client.SIsMember(utils.Ctx, key, UserID).Val()
}

// 查询用户点赞评论状态
func (ls *LikeService) GetCommentLikeStatus(UserID, CommentID string) bool {
	key := fmt.Sprintf("comment:%s:like_users", CommentID)
	return utils.Client.SIsMember(utils.Ctx, key, UserID).Val()
}

// 查询用户点赞回复状态
func (ls *LikeService) GetReplyLikeStatus(UserID, ReplyID string) bool {
	key := fmt.Sprintf("reply:%s:like_users", ReplyID)
	return utils.Client.SIsMember(utils.Ctx, key, UserID).Val()
}

// 查询用户关注状态
func (ls *LikeService) GetFollowStatus(UserID, FolloweeID string) bool {
	key := fmt.Sprintf("User:%s:followers", FolloweeID)
	return utils.Client.SIsMember(utils.Ctx, key, UserID).Val()
}
