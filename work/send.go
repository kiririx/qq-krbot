package work

import (
	"log"
	"qq-krbot/dao"
)

type SendWorker struct {
}

func (*SendWorker) Send() {
	// query to subscribe and user list
	tagUsers, err := dao.SubscribeUserDao.QueryTagAndUser()
	if err != nil {
		log.Println(err)
		return
	}
	for _, tu := range tagUsers {

	}
}
