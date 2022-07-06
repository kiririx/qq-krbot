package work

import (
	"log"
	"qq-krbot/dao"
	"qq-krbot/handler"
	"time"
)

type Collect struct {
}

// Start 采集pixiv图片
func (*Collect) Start() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * 1):
				tags, err := dao.SubscribeDao.Tags()
				if err != nil {
					// 发送到qq错误
					log.Println("采集pixiv记录失败", err)
					continue
				}
				for _, tag := range tags {
					go func(_tag string) {
						i := 0
						for {
							select {
							case <-time.After(time.Second * 1):
								urls, err := handler.GetImgUrlForSearch(_tag, i*handler.PixivPageSize)
								if err != nil {
									// 发送到qq错误
									log.Println("采集pixiv记录失败", err)
									continue
								}
								for _, url := range urls {
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
		}
	}()
}
