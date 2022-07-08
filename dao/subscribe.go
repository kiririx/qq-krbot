package dao

var SubscribeDao = &_SubscribeDao{}
var SubscribeUserDao = &_SubscribeUserDao{}

type _SubscribeDao struct{}
type _SubscribeUserDao struct{}

// Save 保存
func (*_SubscribeDao) Save(tag string, qqAccount string) (SubscribeSubject, error) {
	// todo 加行锁防止并发出现qqAccount错误
	tx := Transaction()
	defer tx.Terminate()
	s := SubscribeSubject{
		Tag: tag,
	}
	err := tx.Sql.Save(&s).Error
	if err != nil {
		return SubscribeSubject{}, err
	}
	err = tx.Sql.Save(&SubscribeUser{
		QQAccount: qqAccount,
		SubId:     s.ID,
	}).Error
	if err != nil {
		return SubscribeSubject{}, err
	}
	tx.Commit()
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

type TagAndUser struct {
	Tag       string `gorm:"column:tag"`
	QQAccount string `gorm:"column:qq_account"`
}

func (*_SubscribeUserDao) QueryTagAndUser() ([]TagAndUser, error) {
	result := make([]TagAndUser, 0)
	err := Sql.Raw("select su.qq_account, ss.tag from subscribe_user su, subscribe_subject ss where su.sub_id = ss.id").Scan(&result).Error
	return result, err
}
