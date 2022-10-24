package dao

import "github.com/kiririx/krutils/mapx"

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
	tx.Sql.Where("tag = ?", tag).Take(&s)
	err := tx.Sql.Save(&s).Error
	if err != nil {
		return SubscribeSubject{}, err
	}
	su := SubscribeUser{
		QQAccount: qqAccount,
		SubId:     s.ID,
	}
	tx.Sql.Where("sub_id = ? and qq_account = ?", s.ID, qqAccount).Take(&su)
	err = tx.Sql.Save(&su).Error
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

func (*_SubscribeUserDao) ClearByUser(qqAccount string) (err error) {
	err = Sql.Unscoped().Where("qq_account = ?", qqAccount).Delete(&SubscribeUser{}).Error
	if err != nil {
		return
	}
	// delete from subscribe_subject where id  not in (select sub_id from subscribe_user)
	err = Sql.Raw("delete from subscribe_subject where id not in (select sub_id from subscribe_user)").Scan(nil).Error
	return
}

type TagAndUser struct {
	Tag       string   `gorm:"column:tag"`
	QQAccount string   `gorm:"column:qq_account"`
	UserTag   []string `gorm:"->"`
}

func (*_SubscribeUserDao) QueryTagAndUser() ([]TagAndUser, error) {
	result := make([]TagAndUser, 0)
	dbTmp := make([]TagAndUser, 0)
	err := Sql.Raw("select su.qq_account, ss.tag from subscribe_user su, subscribe_subject ss where su.sub_id = ss.id").Scan(&dbTmp).Error
	if err != nil {
		return nil, err
	}
	taM := make(map[string][]string)
	for _, v := range dbTmp {
		if mapx.ContainsKey(taM, v.QQAccount) {
			taM[v.QQAccount] = func() []string {
				tmp := taM[v.QQAccount]
				tmp = append(tmp, v.Tag)
				return tmp
			}()
		} else {
			taM[v.QQAccount] = []string{v.Tag}
		}
	}
	for k, v := range taM {
		result = append(result, TagAndUser{
			QQAccount: k,
			UserTag:   v,
		})
	}
	return result, err
}
