package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/map_util"
	"github.com/kiririx/krutils/str_util"
	"log"
	"net/http"
	"qq-krbot/handler"
	"time"
)

var (
	LangReflect map[string]int
	BaseReflect map[string]string
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
			// 判断是否是@消息
			if param.CQ == "at" {
				var msg string
				if str_util.Contains(param.Message, "ヘルプ", "帮助", "help", "?", "？") {
					msg = "🌸リカちゃんの使う方🌸\n" +
						"1. 翻訳：『日中　おはようございます』" + "\n" +
						"2. エロな絵（開発中）" + "\n" +
						"3. 天気予報（開発中）" + "\n" +
						"4. AIリプ（開発中）" + "\n" +
						"5. DNFゴールド：『ゴールド』『金币』『比例』を送信する" + "\n" +
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
	_, _ = http_util.Client(time.Second*4).PostJSON(url, map[string]any{
		"group_id": str_util.ToStr(groupId),
		"message":  msg,
	})
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
