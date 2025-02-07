package dao

import (
	"KnowEase/models"

	"gorm.io/gorm"
)

func ProvideEmailDao(db *gorm.DB) EmailDaoInterface {
	return NewEmailDao(db)
}
func ProvideLikeDao(db *gorm.DB) LikeDaoInterface {
	return NewLikeDao(db)
}
func ProvidePostDao(db *gorm.DB) PostDaoInterface {
	return NewPostDao(db)
}
func ProvideUserDao(db *gorm.DB) UserDaoInterface {
	return NewUserDao(db)
}

type LikeDaoInterface interface {
	GetLikeRecord(Userid string) ([]models.UserLikeHistory, error)
	GetSaveRecord(Userid string) ([]models.UserSaveHistory, error)
	GetViewRecord(Userid string) ([]models.UserViewHistory, error)
	SyncPostLikeToDB(PostID string, LikeCount int) error
	SyncPostCollectToDB(PostID string, SaveCount int) error
	SyncCommentLikeToDB(CommentID string, LikeCount int) error
	SyncReplyLikeToDB(ReplyID string, LikeCount int) error
	SyncPostViewToDB(PostID string, ViewCount int) error
	SyncUserLikesToDB(UserID string, LikeCount int) error
	SyncPostSaveToDB(PostID string, SaveCount int) error
	SyncUserLikeHistoryToDB(Record *models.UserLikeHistory) error
	DeleteUserLikeHistory(PostID string) error
	DeleteUserViewHistory(PostID string) error
	SyncUserSaveHistoryToDB(Record *models.UserSaveHistory) error
	SyncUserViewHistoryToDB(Record *models.UserViewHistory) error
	DeleteUserSaveHistory(PostID string) error
	SearchUserFollowee(UserID string) ([]string, error)
	SearchUserFollower(UserID string) ([]string, error)
	SearchPostmessageFromPost(PostID string) (*models.PostMessage, error)
	SearchCommentmessageFromComment(CommentID string) (*models.Comment, error)
	SyncMessageToDB(Body *models.Message) error
	SearchUnreadMessage(UserID, Tag string) ([]models.Message, error)
	UpdateMessageStatus(UserID, Tag string) error
	SyncUserFollowersToDB(UserID string, FollowerCount int) error
	SyncUserFolloweesToDB(UserID string, FolloweeCount int) error
	SyncFollowMessageToDB(FollowMessage *models.FollowMessage) error
	DeleteFollowMessage(UserID string) error
	SearchCommentByID(CommentID string) (string, string, error)
	SearchReplyByID(ReplyID string) (string, string, error)
}
type PostDaoInterface interface {
	SyncPostBodyToDB(Body *models.PostMessage) error
	DeletePostBody(PostID []string) []error
	SyncCommentBodyToDB(Body *models.Comment) error
	DeleteComment(CommentID string) error
	DeleteReply(ReplyID string) error
	SearchPostID(PostID string) error
	SearchCommentID(CommentID string) error
	SyncReplyToDB(Body *models.Reply) error
	SearchReplyID(ReplyID string) error
	DeleteAllComment(PostID []string) []error
	DeleteAllReply(CommentID string) error
	SearchAllPost() ([]models.PostMessage, error)
	SearchAllComments() ([]models.Comment, error)
	SearchALLReplys() ([]models.Reply, error)
	SearchPostByID(PostID string) (*models.PostMessage, error)
	SearchViewRecord(UserID string) ([]models.UserViewHistory, error)
	SearchUnviewedPost(UserID string) ([]models.PostMessage, error)
	SearchUnviewedPostByTag(UserID, Tag string) ([]models.PostMessage, error)
	SearchCountOfTag(PostID []string) (map[string]int, []error)
	SearchAllComment(PostID string) ([]models.Comment, error)
	SearchALLReply(CommentID string) ([]models.Reply, error)
	SearchCommentByID(CommentID string) (models.Comment, error)
}
type UserDaoInterface interface {
	GetUserFromEmail(email string) (*models.User, error)
	GetUserFromID(UserID string) (*models.Usermessage, error)
	CreateNewUser(user *models.User) error
	ChangePassword(user *models.User) error
	SearchUserid(UserID string) error
	ChangeUserBackground(UserID, Newbackground string) error
	ChangeUserPicture(UserID, NewPicture string) error
	ChangeUserPassword(UserID, NewPassword string) error
	ChangeUserEmail(UserID, NewEmail string) error
	ChangeUsername(UserID, NewName string) error
	SearchAllUser() ([]string, error)
}
type EmailDaoInterface interface {
	WriteCode(email, code string) error
	SearchVerificationCode(code string) (*models.Emailverify, error)
}
