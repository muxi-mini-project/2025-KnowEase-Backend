package dao

import (
	"KnowEase/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PostDao struct {
	db *gorm.DB
}

func NewPostDao(db *gorm.DB) *PostDao {
	return &PostDao{db: db}
}

// 将帖子内容写入数据库
func (pd *PostDao) SyncPostBodyToDB(Body *models.PostMessage) error {
	Body.CreateAt = time.Now()
	r := pd.db.Create(Body)
	return r.Error
}

// 将帖子内容从数据库中删除
func (pd *PostDao) DeletePostBody(PostID []string) []error {
	var Err []error
	for _, id := range PostID {
		err := pd.db.Delete(&models.PostMessage{}, "post_id = ?", id).Error
		if err != nil {
			Err = append(Err, err)
		}
	}
	return Err
}

// 将评论信息写入数据库
func (pd *PostDao) SyncCommentBodyToDB(Body *models.Comment) error {
	Body.CreatAt = time.Now()
	r := pd.db.Create(Body)
	return r.Error
}

// 将评论删除
func (pd *PostDao) DeleteComment(CommentID string) error {
	err := pd.db.Delete(&models.Comment{}, "comment_id = ?", CommentID).Error
	return err
}

// 将回复删除
func (pd *PostDao) DeleteReply(ReplyID string) error {
	err := pd.db.Delete(&models.Reply{}, "reply_id = ?", ReplyID).Error
	return err
}

// 查找帖子id
func (pd *PostDao) SearchPostID(PostID string) error {
	return pd.db.Where(&models.PostMessage{}, "post_id = ?", PostID).Error
}

// 查找评论id
func (pd *PostDao) SearchCommentID(CommentID string) error {
	return pd.db.Where(&models.Comment{}, "comment_id = ?", CommentID).Error
}

// 将回复内容写入数据库
func (pd *PostDao) SyncReplyToDB(Body *models.Reply) error {
	Body.CreatAt = time.Now()
	r := pd.db.Create(Body)
	return r.Error
}

// 查找回复id
func (pd *PostDao) SearchReplyID(ReplyID string) error {
	return pd.db.Where(&models.Reply{}, "reply_id = ?", ReplyID).Error
}

// 删除帖子所有评论
func (pd *PostDao) DeleteAllComment(PostID []string) []error {
	var Err []error
	for _, id := range PostID {
		err := pd.db.Delete(&models.Comment{}, "post_id = ?", id).Error
		if err != nil {
			Err = append(Err, err)
		}
	}
	return Err
}

// 删除评论所有回复
func (pd *PostDao) DeleteAllReply(CommentID string) error {
	err := pd.db.Delete(&models.Reply{}, "comment_id = ?", CommentID).Error
	return err
}

// 查询所有帖子
func (pd *PostDao) SearchAllPost() ([]models.PostMessage, error) {
	var Posts []models.PostMessage
	err := pd.db.Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

// 查询所有评论
func (pd *PostDao) SearchAllComments() ([]models.Comment, error) {
	var Posts []models.Comment
	err := pd.db.Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

// 查询所有评论
func (pd *PostDao) SearchALLReplys() ([]models.Reply, error) {
	var Posts []models.Reply
	err := pd.db.Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

// 根据id查找帖子
func (pd *PostDao) SearchPostByID(PostID string) (*models.PostMessage, error) {
	var Post models.PostMessage
	err := pd.db.Where("post_id = ?", PostID).First(&Post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &Post, nil

}

// 查询用户浏览记录
func (pd *PostDao) SearchViewRecord(UserID string) ([]models.UserViewHistory, error) {
	var Record []models.UserViewHistory
	err := pd.db.Where("user_id = ?", UserID).Find(&Record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return Record, nil
}

// 查询所有用户未浏览的帖子
func (pd *PostDao) SearchUnviewedPost(UserID string) ([]models.PostMessage, error) {
	Record, err := pd.SearchViewRecord(UserID)
	if err != nil {
		return nil, err
	}
	var RecordID []string
	for _, Post := range Record {
		RecordID = append(RecordID, Post.PostID)
	}
	var Posts []models.PostMessage
	if err := pd.db.Where("post_id NOT IN ?", RecordID).Find(&Posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find unviewes posts")
	}
	return Posts, nil
}

// 查询某一tag的未浏览帖子
func (pd *PostDao) SearchUnviewedPostByTag(UserID, Tag string) ([]models.PostMessage, error) {
	Record, err := pd.SearchViewRecord(UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to search user view history")
	}
	var RecordID []string
	for _, Post := range Record {
		RecordID = append(RecordID, Post.PostID)
	}
	var Posts []models.PostMessage
	if err := pd.db.Where("post_id NOT IN ? AND tag = ?", RecordID, Tag).Find(&Posts); err != nil {
		return nil, fmt.Errorf("failed to find unviewed post by tag")
	}
	return Posts, nil
}

// 查询所有包含某个tag的数据
func (pd *PostDao) SearchCountOfTag(PostID []string) (map[string]int, []error) {
	var Err []error
	var Count int64
	tags := make(map[string]int)
	Tag := [4]string{"校园", "生活", "美食", "绘画"}
	for i := 0; i < 4; i++ {
		if err := pd.db.Model(&models.PostMessage{}).Where("post_id IN ? AND tag= ?", PostID, Tag[i]).Count(&Count).Error; err != nil {
			Err = append(Err, err)
		}
		tags[Tag[i]] = int(Count)
	}
	return tags, Err
}

// 查询帖子的所有评论信息
func (pd *PostDao) SearchAllComment(PostID string) ([]models.Comment, error) {
	var Comments []models.Comment
	err := pd.db.Where("post_id = ?", PostID).Find(&Comments).Error
	return Comments, err
}

// 查询评论的所有回复
func (pd *PostDao) SearchALLReply(CommentID string) ([]models.Reply, error) {
	var Replys []models.Reply
	err := pd.db.Where("comment_id = ?", CommentID).Find(&Replys).Error
	return Replys, err
}

// 查询评论
func (pd *PostDao) SearchCommentByID(CommentID string) (models.Comment, error) {
	var Comment models.Comment
	err := pd.db.Where("comment_id = ?", CommentID).Find(&Comment).Error
	return Comment, err
}
