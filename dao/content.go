package dao

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ContentDao = &_ContentDao{}
)

type _ContentDao struct {
}

func (*_ContentDao) QueryRandAndLimit(tag string, limit int) ([]ContentModel, error) {
	var contentModels []ContentModel
	err := Sql.Order("RAND()").Find(&contentModels).Limit(limit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contentModels, nil
		}
		return contentModels, err
	}
	return contentModels, nil
}
