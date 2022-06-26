package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/map_util"
	"github.com/kiririx/krutils/str_util"
	"io/ioutil"
	"log"
	"net/http"
	"qq-krbot/handler"
	"time"
)

var (
	LangReflect map[string]int
	BaseReflect map[string]string
	AnimeImages []string
)

func init() {
	// 语言翻译映射
	LangReflect = make(map[string]int)
	LangReflect["中日"] = cj
	LangReflect["中英"] = ce
	LangReflect["英中"] = ec
	LangReflect["日中"] = jc
	LangReflect["日英"] = je
	LangReflect["英日"] = ej
	// 基本句子映射
	BaseReflect = make(map[string]string)
	BaseReflect["おはよ"] = "おはありです！(｡•ᴗ-)_☀️"
	BaseReflect["こんばんは"] = "こんばんは！(◍•ᴗ•◍)ﾉ"
	BaseReflect["ありがとう"] = "どういたしまして🌷₍ᐢ •͈ ༝ •͈ ᐢ₎♡️"
	BaseReflect["おやすみ"] = "おやすみなさいね〜🌙️"
	BaseReflect["早上好"] = "早上好哦🌷️(◍•ᴗ•◍)ﾉ"
	BaseReflect["中午好"] = "中午好哦🌷️(◍•ᴗ•◍)ﾉ"
	BaseReflect["晚上好"] = "晚上好哦🌷️(◍•ᴗ•◍)ﾉ"
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

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/bot", func(c *gin.Context) {
		param := &Param{}
		err := c.ShouldBindJSON(param)
		param.Parse()
		if err != nil {
			return
		}
		if param.PostType == "meta_event" {
			fmt.Println(param.Message)
			return
		}
		if param.PostType == "message" {
			log.Println("接收消息:", param.Message)
			log.Println("接收消息:", param.RawMessage)
			// 判断是否是翻译消息
			prefix := str_util.SubStr(param.Message, -1, 2)
			if map_util.ContainsKey(LangReflect, prefix) {
				trans(param.Message, param.GroupId, LangReflect[prefix])
				return
			}
			if str_util.Prefix(param.Message, "#") {
				tag := str_util.SubStr(param.Message, 1, -1)
				m, err := handler.Search(tag)
				if err != nil {
					fmt.Println(err.Error())
					return
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
					Error(err, param.GroupId)
					return
				}
				sendMsg := fmt.Sprintf("[CQ:image,file=http://127.0.0.1:10011/photo/%v,subType=0]", imgName)
				sendToGroup(param.GroupId, sendMsg)
				return
			}
			// 判断是否是@消息
			if param.CQ == "at" {
				var msg string
				if str_util.Contains(param.Message, "ヘルプ", "帮助", "help", "?", "？") {
					msg = "🌸リカちゃんの使う方🌸\n" +
						"1. 翻訳：『日中　おはようございます』" + "\n" +
						"2. エロな絵: 『# 可可萝』（# searchKey）を送信する" + "\n" +
						"3. 天気予報（開発中）" + "\n" +
						"4. AIリプ（開発中）" + "\n" +
						"5. DNFゴールド:『ゴールド』『金币』『比例』を送信する" + "\n" +
						"6. ......" + "\n"
				} else if str_util.Contains(param.Message, "DNF", "dnf", "ゴールド", "金币", "比例") {
					glod, err := handler.DNFHandler().Gold()
					if err != nil {
						Error(err, param.GroupId)
						return
					}
					msg = "🌸跨六(uu898): 🌸\n\n"
					for _, v := range glod {
						msg += "比例: " + str_util.ToStr(v.Scale) + "/1    形式: " + v.TradeType + "\n"
					}
				} else if str_util.ContainsSlice(param.Message, map_util.Keys(BaseReflect)) {
					keys := map_util.GetContainedKeys(param.Message, BaseReflect)
					if len(keys) > 0 {
						msg = BaseReflect[keys[0]]
					}
				} else if str_util.Contains(param.Message, "色图", "涩图") {
					fileName := AnimeImages[algo_util.RandomInt(0, len(AnimeImages))]
					sendMsg := fmt.Sprintf("[CQ:image,file=http://127.0.0.1:10011/photo/%v,subType=0]", fileName)
					sendToGroup(param.GroupId, sendMsg)
					return
				} else {
					msg = "何か用事がありますか。"
				}
				sendToGroup(param.GroupId, fmt.Sprintf("[CQ:at,qq=%v] %v", param.UserId, msg))
			}
			// 判断是否是私聊消息
			if param.MessageType == "private" {
				// 判断是否是@消息
				fmt.Println("===>" + param.SubType)
			}
		}
	})
	r.Use(static.Serve("/photo", static.LocalFile("./photo", true)))

	r.Run(":10011") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func trans(text string, groupId int64, transType int) {
	msg := str_util.SubStr(text, 3, -1)
	// translate(msg)
	result := "🌸リカちゃんの翻訳結果🌸：" + Translate(transType, msg)
	// send
	sendToGroup(groupId, result)
}

func sendToGroup(groupId int64, msg string) {
	url := CqHttp + "/send_group_msg"
	resp, err := http_util.Client(time.Second*10).PostJSON(url, map[string]any{
		"group_id": str_util.ToStr(groupId),
		"message":  msg,
	})
	if err != nil {
		Error(err, groupId)
		return
	}
	fmt.Println(resp)
}

func atToGroup(groupId, userId int64, msg string) {
	url := CqHttp + "/send_msg"
	_, _ = http_util.Client(time.Second*4).PostJSON(url, map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"message":  msg,
	})
}

func Error(err error, groupId int64) {
	if err != nil {
		sendToGroup(groupId, "🌸リカちゃんが壊れた：error => 🌸"+err.Error())
	}
}
