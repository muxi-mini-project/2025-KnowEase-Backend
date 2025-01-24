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
	err := ld.db.Model(&models.UserLikes{}).Where("user_id = ?", UserID).Update("like_count", LikeCount).Error
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
func (ld *LikeDao) SearchUnreadMessage(UserID string) ([]models.Message, error) {
	var Messages []models.Message
	err := ld.db.Model(&models.Message{}).Where("user_id = ? AND status = ?", UserID, "unread").Find(&Messages).Error
	return Messages, err
}

// 更新消息状态
func (ld *LikeDao) UpdateMessageStatus(UserID string) error {
	err := ld.db.Model(&models.Message{}).Where("user_id = ?", UserID).Update("status", "read").Error
	return err
}
