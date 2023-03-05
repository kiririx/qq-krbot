package trigger

import (
	"qq-krbot/req"
	"qq-krbot/trigger/resp"
)

var (
	Triggers []Trigger
)

const at = "at"         // 群组里@某人的触发器
const pr = "pr"         // 非群组，单人私聊时的触发器
const Master = "master" // 机器人主人身份，处理一些特殊数据，如密码管理

type Trigger struct {
	Cq        string
	Condition func(*req.Param) bool
	// matchType string
	// word      []string
	Callback func(*req.Param) (string, error)
}

func addTrigger(cq string, condition func(*req.Param) bool, callback func(*req.Param) (string, error)) {
	Triggers = append(Triggers, Trigger{
		Cq:        cq,
		Condition: condition,
		Callback:  callback,
	})
}

func init() {
	addTrigger(pr, Help, resp.Help)                       // 帮助
	addTrigger(pr, DNFGold, resp.DNFGold)                 // dnf金币
	addTrigger(pr, Translate, resp.Translate)             // 翻译
	addTrigger(pr, EroImagesSearch, resp.EroImagesSearch) // 色图搜索
	addTrigger(pr, EroImages, resp.EroImages)
	addTrigger(pr, Text, resp.Text)                         //
	addTrigger(pr, SubscribePixiv, resp.SubscribePixiv)     // 订阅pixiv
	addTrigger(pr, UnSubscribePixiv, resp.UnSubscribePixiv) // 取消订阅pixiv
	addTrigger(pr, Health, resp.Health)                     // ping
	addTrigger(Master, PasswdManage, resp.PasswdManage)     // 密码管理
}
