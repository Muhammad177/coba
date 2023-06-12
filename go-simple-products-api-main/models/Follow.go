package models

type Follow struct {
	ID     uint `gorm:"primary_key"`
	UserID int  `json:"user_id" form:"user_id"`
	User   User `json:"user"`
}
