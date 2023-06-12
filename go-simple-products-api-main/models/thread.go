package models

import (
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model
	Title          string    `json:"title" form:"title"`
	Topic          string    `json:"topic" form:"topic"`
	Content        string    `json:"content" form:"content"`
	File           string    `json:"file" form:"file"`
	UserID         int       `json:"user_id" form:"user_id"`
	User           User      `json:"user"`
	Comments       []Comment `json:"comment"`
	BookmarkedUser []User    `gorm:"many2many: thread_user_assoc"`
}
type AllThread struct {
	ID      uint    `gorm:"primary_key"`
	Title   string  `json:"title" form:"title"`
	Topic   string  `json:"topic" form:"topic"`
	Content string  `json:"content" form:"content"`
	File    string  `json:"file" form:"file"`
	UserID  int     `json:"user_id" form:"user_id"`
	User    AllUser `json:"user"`
}
type ThreadUser struct {
	ID    uint   `gorm:"primary_key"`
	Title string `json:"title" form:"title"`
	Topic string `json:"topic" form:"topic"`
}

func ConvertThreadUser(thread *Thread) ThreadUser {
	return ThreadUser{
		ID:    thread.ID,
		Title: thread.Title,
		Topic: thread.Topic,
	}
}

type ThreadResponse struct {
	ID       uint              `gorm:"primary_key"`
	Title    string            `json:"title" form:"title"`
	Topic    string            `json:"topic" form:"topic"`
	Comments []CommentResponse `json:"comments"`
}

func ConvertThreadToThreadResponse(thread *Thread) ThreadResponse {
	comments := make([]CommentResponse, len(thread.Comments))
	for i, comment := range thread.Comments {
		comments[i] = ConvertCommentToCommentResponse(&comment)
	}

	return ThreadResponse{
		ID:       thread.ID,
		Title:    thread.Title,
		Topic:    thread.Topic,
		Comments: comments,
	}
}
func ConverThreadToAllThread(thread *Thread) AllThread {

	return AllThread{
		ID:      thread.ID,
		Title:   thread.Title,
		Topic:   thread.Topic,
		Content: thread.Content,
		File:    thread.File,
		UserID:  thread.UserID,
		User:    ConvertUserToAllUser(&thread.User),
	}
}
