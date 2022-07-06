package dao

import "github.com/kiririx/krutils/str_util"

var SubscribeDao = &_SubscribeDao{}

type _SubscribeDao struct{}

// Save 保存
func (*_SubscribeDao) Save(tag string, qqAccount string) (SubscribeSubject, error) {
	// todo 加行锁防止并发出现qqAccount错误
	s := SubscribeSubject{}
	Sql.Where("tag = ?", tag).Take(&s)
	if s.ID > 0 && !str_util.Contains(s.QQAccount, qqAccount) {
		qqAccount += "," + qqAccount
	}
	s.Tag = tag
	s.QQAccount = qqAccount
	s.Active = true
	err := Sql.Save(&s).Error
	return s, err
}

func (*_SubscribeDao) Tags() ([]string, error) {
	subs := make([]SubscribeSubject, 0)
	err := Sql.Find(&subs).Error
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0, len(subs))
	for _, s := range subs {
		tags = append(tags, s.Tag)
	}
	return tags, nil
}
