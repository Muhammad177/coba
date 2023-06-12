package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Comment  string `json:"comment" form:"comment"`
	ThreadID int    `json:"thread_id" form:"thread_id"`
	Thread   Thread `json:"thread"`
	UserID   int    `json:"user_id" form:"user_id"`
	User     User   `json:"user"`
}
type CommentResponse struct {
	Comment  string `json:"comment" form:"comment"`
	ThreadID int    `json:"thread_id" form:"thread_id"`
	UserID   int    `json:"user_id" form:"user_id"`
}

func ConvertCommentToCommentResponse(comment *Comment) CommentResponse {
	return CommentResponse{
		Comment:  comment.Comment,
		ThreadID: comment.ThreadID,
		UserID:   comment.UserID,
	}
}

