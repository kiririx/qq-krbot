package work

import (
	"github.com/kiririx/krutils/algo_util"
	"log"
	"qq-krbot/dao"
	"qq-krbot/env"
	"qq-krbot/qqutil"
	"time"
)

type SendWorker struct {
}

func (*SendWorker) Start() {
	go func() {
		// query to subscribe and user list
		tagUsers, err := dao.SubscribeUserDao.QueryTagAndUser()
		if err != nil {
			log.Println(err)
			return
		}
		for {
			time.Sleep(time.Second * 5)
			tu := tagUsers[algo_util.RandomInt(0, len(tagUsers)-1)]
			qqutil.SendPrivateMessage(tu.QQAccount, qqutil.QQMsg{
				CQ: "image",
				FileURL: func() string {
					l := len(dao.FileList[tu.Tag])
					path := dao.FileList[tu.Tag][algo_util.RandomInt(0, l-1)]
					return "http://127.0.0.1:" + env.Conf["serve.port"] + "/" + path
				}(),
			})
		}
	}()

}
