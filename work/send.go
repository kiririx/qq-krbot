package work

import (
	"github.com/kiririx/krutils/algox"
	"github.com/kiririx/krutils/convertx"
	"qq-krbot/dao"
	"qq-krbot/env"
	"qq-krbot/qqutil"
	"time"
)

type SendWorker struct {
}

func (*SendWorker) Start() {
	go func() {
		update := false
		// query to subscribe and user list
		tagUsers, _ := dao.SubscribeUserDao.QueryTagAndUser()
		upTicker := time.NewTicker(time.Minute)
		sdTicker := time.NewTicker(time.Second * time.Duration(convertx.ToInt(env.Conf["send.time"])))
		go func() {
			for {
				<-upTicker.C
				update = true
			}
		}()
		for {
			<-sdTicker.C
			// ctx, cancelFunc := context.WithCancel(context.TODO())
			// cancelFunc()
			if update {
				tagUsers, _ = dao.SubscribeUserDao.QueryTagAndUser()
				update = false
			}
			if len(tagUsers) < 1 {
				continue
			}
			for _, tus := range tagUsers {
				go func(qqAccount string, tags []string) {
					qqutil.SendPrivateMessage(qqAccount, qqutil.QQMsg{
						CQ: "image",
						FileURL: func() string {
							tagMap := dao.FileList.Get(tags[algox.RandomInt(0, len(tags)-1)])
							if tagMap != nil {
								path := tagMap.Get(algox.RandomInt(0, tagMap.Len()-1))
								if path == "" {
									return ""
								}
								return "http://127.0.0.1:" + env.Conf["serve.port"] + "/" + path
							}
							return ""
						}(),
					})
				}(tus.QQAccount, tus.UserTag)
			}
		}
	}()

}
