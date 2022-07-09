package qqutil

import (
	"fmt"
	"github.com/kiririx/krutils/http_util"
	"log"
	"qq-krbot/env"
	"time"
)

type QQMsg struct {
	CQ      string
	FileURL string
	Message string
}

func SendPrivateMessage(targetQQ string, msg QQMsg) {
	cqHttp := env.Conf["cqhttp.url"]
	if cqHttp == "" {
		panic("CQHttp地址未配置")
	}
	if msg.CQ == "image" && msg.FileURL == "" {
		return
	}
	url := cqHttp + "/send_private_msg"
	resp, err := http_util.Client().Timeout(time.Second*30).Headers(map[string]string{
		"content-type": "application/json",
	}).PostString(url, map[string]any{
		"message": func() string {
			if msg.CQ == "image" {
				return fmt.Sprintf("[CQ:image,file=%v,subType=%v]", msg.FileURL, 0)
			}
			return msg.Message
		}(),
		"user_id":     targetQQ,
		"auto_escape": false,
	})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("cq-http-resp => ", resp)
}
