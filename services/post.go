package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type PostService struct {
	PostDao *dao.PostDao
	LikeDao *dao.LikeDao
	UserDao *dao.UserDao
}

func NewPostService(PostDao *dao.PostDao, LikeDao *dao.LikeDao, UserDao *dao.UserDao) *PostService {
	return &PostService{PostDao: PostDao, LikeDao: LikeDao, UserDao: UserDao}
}

// 发布帖子
func (ps *PostService) PublishPost(Post models.PostMessage) error {
	return ps.PostDao.SyncPostBodyToDB(&Post)
}

// 发布评论
func (ps *PostService) PublishComment(Comment models.Comment) error {
	return ps.PostDao.SyncCommentBodyToDB(&Comment)
}

// 发布评论
func (ps *PostService) PublishReply(Reply models.Reply) error {
	return ps.PostDao.SyncReplyToDB(&Reply)
}

// 删除评论
func (ps *PostService) DeleteComment(CommentID string) error {
	return ps.PostDao.DeleteComment(CommentID)
}

// 删除回复
func (ps *PostService) DeleteReply(ReplyID string) error {
	return ps.PostDao.DeleteReply(ReplyID)
}

// 删除帖子
func (ps *PostService) DeletePost(PostID []string) []error {
	return ps.PostDao.DeletePostBody(PostID)
}

// 删除帖子相关评论
func (ps *PostService) DeletePostComment(PostID []string) []error {
	return ps.PostDao.DeleteAllComment(PostID)
}

// 查找相关帖子信息-根据id
func (ps *PostService) SearchPostByID(ID []string) ([]models.PostMessage, []error) {
	var PostMessages []models.PostMessage
	var Err []error
	for _, PostID := range ID {
		PostMessage, err := ps.PostDao.SearchPostByID(PostID)
		if err != nil {
			Err = append(Err, err)
		}
		PostMessages = append(PostMessages, *PostMessage)
	}
	if PostMessages == nil {
		return nil, Err
	}
	if Err != nil {
		log.Println("something went wrong while querying recommended posts:%w", Err)
	}
	return PostMessages, nil
}

// 推荐的加权
func (ps *PostService) WeightedRecommendation(UserID string) ([]models.PostRecommendLevel, error) {
	Record, err := ps.PostDao.SearchViewRecord(UserID)
	if err != nil {
		return nil, err
	}
	var PostIDS []string
	for _, Post := range Record {
		PostIDS = append(PostIDS, Post.PostID)
	}
	var WeightedPosts []models.PostRecommendLevel
	tagCount, Err := ps.PostDao.SearchCountOfTag(PostIDS)
	if Err != nil {
		return nil, fmt.Errorf("failed to count tags")
	}
	PostMessage, err := ps.PostDao.SearchUnviewedPost(UserID)
	if err != nil {
		return nil, err
	}
	//计算所有帖子的权重
	for _, Posts := range PostMessage {
		weight := (Posts.LikeCount*3+Posts.SaveCount*4+Posts.ViewCount*3)*4/100 + tagCount[Posts.Tag]*6/10
		WeightedPosts = append(WeightedPosts, models.PostRecommendLevel{
			PostID: Posts.PostID,
			Weight: weight,
		})
	}
	//根据权重排序
	for i := 0; i < len(WeightedPosts); i++ {
		for j := 0; j < len(WeightedPosts)-i; j++ {
			if WeightedPosts[i].Weight < WeightedPosts[j].Weight {
				WeightedPosts[i] = WeightedPosts[j]
			}
		}
	}
	return WeightedPosts, nil
}

// 查询某一tag的未浏览帖子
func (ps *PostService) SearchUnviewedPostsByTag(UserID, Tag string) ([]models.PostMessage, error) {
	return ps.PostDao.SearchUnviewedPostByTag(UserID, Tag)
}
func (ps *PostService) DeleteAllReply(CommentID string) error {
	return ps.PostDao.DeleteAllReply(CommentID)
}

// 查询帖子的所有评论
func (ps *PostService) GetAllComment(PostID string) (*models.PostMessage, error) {
	PostMessage, err := ps.PostDao.SearchPostByID(PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to find this Post")
	}
	//查询所有的评论
	Comments, err := ps.PostDao.SearchAllComment(PostMessage.PostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return PostMessage, nil
		}
		return nil, err
	}
	for i := 0; i < len(Comments); i++ {
		Reply, _ := ps.PostDao.SearchALLReply(Comments[i].CommentID)
		Comments[i].Reply = Reply
		for j := 0; j < len(Reply); j++ {
			Replys, _ := ps.PostDao.SearchALLReply(Reply[j].ReplyID)
			Reply[j].Reply = Replys
		}
	}
	PostMessage.Comment = Comments
	return PostMessage, nil

}
func (pd *PostService) SearchPosterMessage(UserID string) (string, string, error) {
	PosterMessage, err := pd.UserDao.GetUserFromID(UserID)
	if err != nil {
		return "", "", err
	}
	return PosterMessage.Username, PosterMessage.ImageURL, nil
}

// 查询未读消息
func (pd *PostService) SearchAllUnreadMessage(UserID string) ([]models.Message, error) {
	Messages, err := pd.LikeDao.SearchUnreadMessage(UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return Messages, nil
}

// 更新消息状态
func (pd *PostService) UpdateMessageStatus(UserID string) error {
	return pd.LikeDao.UpdateMessageStatus(UserID)
}
