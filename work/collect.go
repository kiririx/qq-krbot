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
	tags, err := dao.SubscribeDao.Tags()
	if err != nil {
		// 发送到qq错误
		log.Println("获取 pixiv tag 失败", err)
		return
	}
	for _, tag := range tags {
		go func(_tag string) {
			i := 0
			for {
				select {
				case <-time.After(time.Second * 30):
					log.Println("当前offset：" + str_util.ToStr(i*handler.PixivPageSize))
					urls, err := handler.GetImgUrlForSearch(_tag, i*handler.PixivPageSize)
					if err != nil {
						// 发送到qq错误
						log.Println("采集pixiv记录失败", err)
						continue
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
				}
				i++
			}
		}(tag)
	}
}
