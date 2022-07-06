package dao

import "gorm.io/gorm"

type ContentModel struct {
	Tag     string `gorm:"column:tag"`
	Content string `gorm:"column:content"`
	gorm.Model
}

func (*ContentModel) TableName() string {
	return "content"
}

// SubscribeSubject 订阅的主题
type SubscribeSubject struct {
	Tag       string `gorm:"column:tag"`
	QQAccount string `gorm:"column:qq_account"`
	Active    bool   `gorm:"column:active"`
	gorm.Model
}

func (*SubscribeSubject) TableName() string {
	return "subscribe_subject"
}
