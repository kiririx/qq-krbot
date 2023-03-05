package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiririx/krutils/httpx"
	"github.com/kiririx/krutils/strx"
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
					qqutil.SendPrivateMessage(strx.ToStr(param.UserId), qqutil.QQMsg{
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
	sendGroupId := strx.ToStr(groupId)
	resp, err := httpx.Client().Timeout(time.Second*30).PostString(url, map[string]any{
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
		qqutil.SendPrivateMessage(strx.ToStr(2187391949), qqutil.QQMsg{
			Message: err.Error(),
			CQ:      "pr",
		})
	}
}
