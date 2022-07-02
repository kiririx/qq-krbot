package dao

import "gorm.io/gorm"

type ContentModel struct {
	Tag     string `gorm:"column:tag"`
	Content string `gorm:"column:content"`
	*gorm.Model
}

func (*ContentModel) TableName() string {
	return "content"
}
