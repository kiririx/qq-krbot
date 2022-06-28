package resp

import (
	"fmt"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/str_util"
	"io/ioutil"
	"qq-krbot/handler"
	"qq-krbot/req"
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
	return "🌸リカちゃんの使う方🌸\n" +
		"1. 翻訳：『日中　おはようございます』" + "\n" +
		"2. エロな絵: 『#可可萝』（#searchKey）を送信する" + "\n" +
		"3. 天気予報（開発中）" + "\n" +
		"4. AIリプ（開発中）" + "\n" +
		"5. DNFゴールド:『ゴールド』『金币』『比例』を送信する" + "\n" +
		"6. ......" + "\n", nil
}

func Translate(r *req.Param) (string, error) {
	transType := str_util.SubStr(r.Message, -1, 2)
	transText := str_util.SubStr(r.Message, 3, -1)
	result := "🌸リカちゃんの翻訳結果🌸：" + handler.Translate(handler.LangReflect[transType], transText)
	if r.CQ == "at" {
		return AtResp(result, r.UserId), nil
	}
	return result, nil
}

func EroImagesSearch(r *req.Param) (string, error) {
	tag := str_util.SubStr(r.Message, 1, -1)
	m, err := handler.Search(tag)
	if err != nil {
		return "", err
	}
	fmt.Println(m)
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
	p := photos[algo_util.RandomInt(0, len(photos))]
	imgName, err := handler.DownloadImg(p)
	if err != nil {
		return "", err
	}
	return ImgResp("http://127.0.0.1:10011/photo/"+imgName, 0), nil
}

func EroImages(*req.Param) (string, error) {
	fileName := AnimeImages[algo_util.RandomInt(0, len(AnimeImages))]
	return ImgResp("http://127.0.0.1:10011/photo/"+fileName, 0), nil
}

func SimpleReflect(*req.Param) (string, error) {
	return "", nil
}
