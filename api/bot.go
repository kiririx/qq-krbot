package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/str_util"
	"log"
	"net/http"
	"qq-krbot/qqutil"
	"qq-krbot/req"
	"qq-krbot/trigger"
	"time"
)

const (
	// CqHttp go-cqhttp 地址
	CqHttp = "http://127.0.0.1:5700"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Bot(c *gin.Context) {
	param := &req.Param{}
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
		for _, tg := range trigger.Triggers {
			if tg.Cq == param.CQ && tg.Condition(param) {
				msg, err := tg.Callback(param)
				if err != nil {
					Error(err, param.GroupId)
					return
				}
				switch tg.Cq {
				case "pr":
					qqutil.SendPrivateMessage(str_util.ToStr(param.UserId), qqutil.QQMsg{
						Message: msg,
						CQ:      tg.Cq,
					})
				case "at":
					sendToGroup(param.GroupId, msg)
				}
				break
			}
		}
	}
}

func sendToGroup(groupId int64, msg string) {
	url := CqHttp + "/send_group_msg"
	sendGroupId := str_util.ToStr(groupId)
	resp, err := http_util.Client().Timeout(time.Second*30).PostString(url, map[string]any{
		"group_id": sendGroupId,
		"message":  msg,
	})
	if err != nil {
		Error(err, groupId)
		return
	}
	fmt.Println("group_id => ", sendGroupId)
	fmt.Println("cq-http-resp => ", resp)

}

func Error(err error, groupId int64) {
	if err != nil {
		log.Println(groupId, "Error => 🌸", err)
	}
}
