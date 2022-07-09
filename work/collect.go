package work

import (
	"github.com/kiririx/krutils/str_util"
	"log"
	"qq-krbot/dao"
	"qq-krbot/handler"
	"time"
)

type CollectWorker struct {
}

// Start 采集pixiv图片
func (*CollectWorker) Start() {
	go func() {
		tagM := dao.SyncMap[string, *int]()
		for {
			upTicker := time.NewTicker(time.Minute * 1)
			coTicker := time.NewTicker(time.Second * 45)
			update := false
			tags, _ := dao.SubscribeDao.Tags()
			go func() {
				for {
					<-upTicker.C
					update = true
				}
			}()
			for {
				<-coTicker.C
				if update {
					tags, _ = dao.SubscribeDao.Tags()
					update = false
				} else {
					for _, tag := range tags {
						if v := tagM.Get(tag); v != nil {
							go func(_tag string) {
								singTagMap := dao.FileList.Get(_tag)
								if singTagMap != nil && singTagMap.Len() > 500 {
									return
								}
								i := 0
								tagM.Put(_tag, &i)
								for {
									log.Println("当前tag: " + _tag + " offset:" + str_util.ToStr(i*handler.PixivPageSize))
									urls, err := handler.GetImgUrlForSearch(_tag, i*handler.PixivPageSize)
									if err != nil {
										// 发送到qq错误
										log.Println("采集pixiv记录失败", err)
										return
									}
									for _, url := range urls {
										time.Sleep(time.Second * 1)
										_, err := handler.DownloadImgAndTag(url, _tag)
										if err != nil {
											// 发送到qq错误
											log.Println("下载pixiv图片失败", err)
											continue
										}
									}
									i++
									if i > 30 {
										tagM.Remove(_tag)
										return
									}
								}
							}(tag)
						}
					}
				}
			}
		}
	}()

}
