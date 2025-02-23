package dao

import (
	"KnowEase/models"
	"time"

	"gorm.io/gorm"
)

type LikeDao struct {
	db *gorm.DB
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{db: db}
}

// 在数据库中查询最近一个月的历史点赞记录
func (ld *LikeDao) GetLikeRecord(Userid string) ([]models.UserLikeHistory, error) {
	oneMonthAgo := time.Now().AddDate(0, 0, -31) // 获取 31天前的时间
	var LikeRecords []models.UserLikeHistory
	if err := ld.db.Order("created_at DESC").Where("user_id = ? AND created_at > ?", Userid, oneMonthAgo).Find(&LikeRecords).Error; err != nil {
		return nil, err
	}
	return LikeRecords, nil
}

// 在数据库中查询历史收藏记录
func (ld *LikeDao) GetSaveRecord(Userid string) ([]models.UserSaveHistory, error) {
	var SaveRecords []models.UserSaveHistory
	if err := ld.db.Order("created_at DESC").Where("user_id = ? ", Userid).Find(&SaveRecords).Error; err != nil {
		return nil, err
	}
	return SaveRecords, nil
}

// 在数据库中查询历史收藏记录
func (ld *LikeDao) GetViewRecord(Userid string) ([]models.UserViewHistory, error) {
	var ViewRecords []models.UserViewHistory
	if err := ld.db.Order("created_at DESC").Where("user_id = ? ", Userid).Find(&ViewRecords).Error; err != nil {
		return nil, err
	}
	return ViewRecords, nil
}

// 将帖子点赞数据写入数据库
func (ld *LikeDao) SyncPostLikeToDB(PostID string, LikeCount int) error {
	err := ld.db.Model(&models.PostMessage{}).Where("post_id = ?", PostID).Update("like_count", LikeCount).Error
	return err
}

// 将帖子收藏数据写入数据库
func (ld *LikeDao) SyncPostCollectToDB(PostID string, SaveCount int) error {
	err := ld.db.Model(&models.PostMessage{}).Where("post_id = ?", PostID).Update("save_count", SaveCount).Error
	return err
}

// 将评论点赞数据写入数据库
func (ld *LikeDao) SyncCommentLikeToDB(CommentID string, LikeCount int) error {
	err := ld.db.Model(&models.Comment{}).Where("comment_id = ?", CommentID).Update("like_count", LikeCount).Error
	return err
}

// 将回复点赞数据写入数据库
func (ld *LikeDao) SyncReplyLikeToDB(ReplyID string, LikeCount int) error {
	err := ld.db.Model(&models.Reply{}).Where("reply_id = ?", ReplyID).Update("like_count", LikeCount).Error
	return err
}

// 将帖子浏览数据写入数据库
func (ld *LikeDao) SyncPostViewToDB(PostID string, ViewCount int) error {
	err := ld.db.Model(&models.PostMessage{}).Where("post_id = ?", PostID).Update("view_count", ViewCount).Error
	return err
}

// 将用户获赞数数据写入数据库
func (ld *LikeDao) SyncUserLikesToDB(UserID string, LikeCount int) error {
	err := ld.db.Model(&models.User{}).Where("id = ?", UserID).Update("like_count", LikeCount).Error
	return err
}

// 将帖子点赞数据写入数据库
func (ld *LikeDao) SyncPostSaveToDB(PostID string, SaveCount int) error {
	err := ld.db.Model(&models.PostMessage{}).Where("post_id = ?", PostID).Update("save_count", SaveCount).Error
	return err
}

// 将用户点赞历史记录写入数据库
func (ld *LikeDao) SyncUserLikeHistoryToDB(Record *models.UserLikeHistory) error {
	Record.CreatedAt = time.Now()
	r := ld.db.Create(Record)
	return r.Error
}

// 将用户点赞历史记录从数据库中删除
func (ld *LikeDao) DeleteUserLikeHistory(PostID string) error {
	err := ld.db.Delete(&models.UserLikeHistory{}, "post_id = ?", PostID).Error
	return err
}

// 将用户浏览历史记录从数据库中删除
func (ld *LikeDao) DeleteUserViewHistory(PostID string) error {
	err := ld.db.Delete(&models.UserViewHistory{}, "post_id = ?", PostID).Error
	return err
}

// 将用户收藏历史记录写入数据库
func (ld *LikeDao) SyncUserSaveHistoryToDB(Record *models.UserSaveHistory) error {
	Record.CreatedAt = time.Now()
	r := ld.db.Create(Record)
	return r.Error
}

// 将用户浏览历史记录写入数据库
func (ld *LikeDao) SyncUserViewHistoryToDB(Record *models.UserViewHistory) error {
	Record.CreateAt = time.Now()
	err := ld.db.Model(&models.UserViewHistory{}).Where("user_id = ? AND post_id = ?", Record.UserID, Record.PostID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r := ld.db.Create(Record)
			return r.Error
		}
		return err
	}
	err = ld.db.Model(&models.UserViewHistory{}).Where("user_id = ? AND post_id = ?", Record.UserID, Record.PostID).Update("create_at", Record.CreateAt).Error
	return err
}

// 将用户收藏历史记录从数据库中删除
func (ld *LikeDao) DeleteUserSaveHistory(PostID string) error {
	err := ld.db.Delete(&models.UserSaveHistory{}, "post_id = ?", PostID).Error
	return err
}

// 根据帖子id查找帖子信息
func (ld *LikeDao) SearchPostmessageFromPost(PostID string) (*models.PostMessage, error) {
	var ID models.PostMessage
	re := ld.db.Where("post_id=?", PostID).Find(&ID)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &ID, nil
}

// 根据评论id查找评论
func (ld *LikeDao) SearchCommentmessageFromComment(CommentID string) (*models.Comment, error) {
	var ID models.Comment
	re := ld.db.Where("comment_id = ?", CommentID).Find(&ID)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &ID, nil
}

// 将消息内容写入数据库
func (ld *LikeDao) SyncMessageToDB(Body *models.Message) error {
	Body.CreateAt = time.Now()
	r := ld.db.Create(Body)
	return r.Error
}

// 查询用户所有未读消息
func (ld *LikeDao) SearchUnreadMessage(UserID,Tag string) ([]models.Message, error) {
	var Messages []models.Message
	err := ld.db.Model(&models.Message{}).Where("user_id = ? AND status = ? AND tag = ?", UserID, "unread",Tag).Find(&Messages).Error
	return Messages, err
}

// 更新消息状态
func (ld *LikeDao) UpdateMessageStatus(UserID,Tag string) error {
	err := ld.db.Model(&models.Message{}).Where("user_id = ? AND tag = ?", UserID,Tag).Update("status", "read").Error
	return err
}

// 将用户粉丝数数据写入数据库
func (ld *LikeDao) SyncUserFollowersToDB(UserID string, FollowerCount int) error {
	err := ld.db.Model(&models.User{}).Where("id = ?", UserID).Update("follower_count", FollowerCount).Error
	return err
}

// 将用户关注数数据写入数据库
func (ld *LikeDao) SyncUserFolloweesToDB(UserID string, FolloweeCount int) error {
	err := ld.db.Model(&models.User{}).Where("id = ?", UserID).Update("followee_count", FolloweeCount).Error
	return err
}

// 将关注信息写入数据库
func (ld *LikeDao) SyncFollowMessageToDB(FollowMessage *models.FollowMessage) error {
	return ld.db.Create(FollowMessage).Error
}

// 将关注信息从数据库里删除
func (ld *LikeDao) DeleteFollowMessage(UserID string) error {
	return ld.db.Delete(&models.FollowMessage{}, "follower_id = ?", UserID).Error
}

// 获取用户关注列表
func (ld *LikeDao) SearchUserFollowee(UserID string) ([]string, error) {
	var FolloweeIDs []string
	err := ld.db.Model(&models.FollowMessage{}).Where("follower_id", UserID).Select("followee_id").Find(&FolloweeIDs).Error
	return FolloweeIDs, err
}

// 获取用户粉丝列表
func (ld *LikeDao) SearchUserFollower(UserID string) ([]string, error) {
	var FollowerIDs []string
	err := ld.db.Model(&models.FollowMessage{}).Where("followee_id", UserID).Select("follower_id").Find(&FollowerIDs).Error
	return FollowerIDs, err
}

// 根据id找评论
func (ld *LikeDao) SearchCommentByID(CommentID string) (string, string, error) {
	var Comment models.Comment
	err := ld.db.Model(&models.Comment{}).Where("comment_id = ?", CommentID).First(&Comment).Error
	if err != nil {
		return "", "", err
	}
	return Comment.CommenterID, Comment.Body, nil
}

// 根据id找回复
func (ld *LikeDao) SearchReplyByID(ReplyID string) (string, string, error) {
	var Reply models.Reply
	err := ld.db.Model(&models.Reply{}).Where("reply_id = ?", ReplyID).First(&Reply).Error
	if err != nil {
		return "", "", err
	}
	return Reply.ReplyerID, Reply.Body, nil
}
