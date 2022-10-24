package work

import (
	"github.com/kiririx/krutils/strx"
	"github.com/kiririx/krutils/sugar"
	"log"
	"qq-krbot/dao"
	"qq-krbot/handler"
	"sync"
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
			coTicker := time.NewTicker(time.Second * 33)
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
				}
				wg := sync.WaitGroup{}
				for _, tag := range tags {
					wg.Add(1)
					go func(_tag string) {
						singTagMap := dao.FileList.Get(_tag)
						if singTagMap != nil && singTagMap.Len() > 500 {
							return
						}
						v := 0
						i := *sugar.Then(tagM.Get(_tag) != nil, tagM.Get(_tag), &v)
						if i > 30 {
							return
						}
						log.Println("当前tag: " + _tag + " offset:" + strx.ToStr(i*handler.PixivPageSize))
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
						tagM.Put(_tag, &i)
						wg.Done()
					}(tag)
				}
				wg.Wait()
			}
		}
	}()

}
