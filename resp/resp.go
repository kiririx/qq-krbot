package resp

import (
	"errors"
	"fmt"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/slice_util"
	"github.com/kiririx/krutils/str_util"
	"io/ioutil"
	"qq-krbot/dao"
	"qq-krbot/handler"
	"qq-krbot/req"
	"time"
)

var (
	AnimeImages []string
)

type Resp struct {
	Message string
	CQ      string
	SubType int
	File    string
}

func AtResp(message string, userId int64) string {
	return fmt.Sprintf("[CQ:at,qq=%v] %v", userId, message)
}

func ImgResp(url string, subType int) string {
	return fmt.Sprintf("[CQ:image,file=%v,subType=%v]", url, subType)
}

func init() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				// 涩图列表
				AnimeImages = make([]string, 0)
				files, _ := ioutil.ReadDir("./photo")
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					AnimeImages = append(AnimeImages, file.Name())
				}
			}
		}
	}()
}

func DNFGold(*req.Param) (string, error) {
	gold, err := handler.DNFHandler().Gold()
	if err != nil {
		return "", err
	}
	msg := "🌸跨六(uu898): 🌸\n\n"
	for _, v := range gold {
		msg += "比例: " + str_util.ToStr(v.Scale) + "/1    形式: " + v.TradeType + "\n"
	}
	return msg, nil
}

func Help(*req.Param) (string, error) {
	return "🌸梨花酱的使用方法🌸\n" +
		"1. 翻译：发送『中日　早上好』" + "\n" +
		"2. 涩图: 发送『#可可萝』（#要搜索的内容）" + "\n" +
		"3. 天气预报（开发中）" + "\n" +
		"4. AI回复（开发中）" + "\n" +
		"5. DNF金币比例: 发送『金币』或『比例』" + "\n" +
		"6. ......" + "\n", nil
}

func Translate(r *req.Param) (string, error) {
	transType := str_util.SubStr(r.Message, -1, 2)
	transText := str_util.SubStr(r.Message, 3, -1)
	result := "🌸翻译结果🌸：" + handler.Translate(handler.LangReflect[transType], transText)
	return AtResp(result, r.UserId), nil
}

func EroImagesSearch(r *req.Param) (string, error) {
	tag := str_util.SubStr(r.Message, 1, -1)
	m, err := handler.Search(tag, handler.PixivPageSize)
	if err != nil {
		return "", err
	}
	photos := make([]string, 0)
	illusts := m["illusts"].([]interface{})
	for _, illust := range illusts {
		image := illust.(map[string]interface{})["meta_single_page"].(map[string]interface{})
		var imgUrl = image["original_image_url"]
		if imgUrl != nil {
			photos = append(photos, imgUrl.(string))
		} else {
			images := illust.(map[string]interface{})["meta_pages"].([]interface{})
			for _, v := range images {
				img := v.(map[string]interface{})
				url := img["image_urls"].(map[string]interface{})["original"]
				if url != nil {
					photos = append(photos, url.(string))
				}
			}
		}
	}
	imgName, err := download(photos)
	if err != nil {
		return "", err
	}
	return ImgResp("http://127.0.0.1:10013/photo/"+imgName, 0), nil
}

func download(photos []string) (string, error) {
	if len(photos) < 1 {
		return "", errors.New("all photos download failed")
	}
	i := algo_util.RandomInt(0, len(photos))
	p := photos[i]
	imgName, err := handler.DownloadImg(p)
	if err != nil {
		imgName, err = download(slice_util.Remove(photos, i))
		if err != nil {
			return "", err
		}
	}
	return imgName, nil
}

func EroImages(*req.Param) (string, error) {
	fileName := AnimeImages[algo_util.RandomInt(0, len(AnimeImages))]
	return ImgResp("http://127.0.0.1:10013/photo/"+fileName, 0), nil
}

func Text(req *req.Param) (string, error) {
	c, err := dao.ContentDao.QueryRandAndLimit("text", 1)
	if err != nil {
		return "", err
	}
	return c[0].Content, nil
}

func Word(req *req.Param) (string, error) {
	c, err := dao.ContentDao.QueryRandAndLimit("word", 10)
	if err != nil {
		return "", err
	}
	contents := ""
	for _, v := range c {
		contents += v.Content + "\n"
	}
	return contents, nil
}

func SimpleReflect(*req.Param) (string, error) {
	return "", nil
}

func SubscribePixiv(param *req.Param) (string, error) {
	tag := str_util.SubStr(param.Message, 3, -1)
	tag = str_util.TrimSpace(tag)
	_, err := dao.SubscribeDao.Save(tag, str_util.ToStr(param.UserId))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("开始订阅「%v」相关的图片，リカ酱会大约1分钟发一次", tag), nil
}

func UnSubscribePixiv(param *req.Param) (string, error) {
	err := dao.SubscribeUserDao.ClearByUser(str_util.ToStr(param.UserId))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("已经取消全部订阅，几分钟后リカ酱不会发送信息，如需重新订阅请输入：订阅 xxx"), nil
}

func Health(param *req.Param) (string, error) {
	return "pong", nil
}
