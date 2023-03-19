package trigger

import (
	"github.com/kiririx/krutils/mapx"
	"github.com/kiririx/krutils/strx"
	"qq-krbot/handler"
	"qq-krbot/req"
)

var (
	BaseReflect map[string]string
)

func init() {

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

func Help(param *req.Param) bool {
	return strx.Equals(param.Message, "ヘルプ", "帮助", "help", "?", "？")
}

func DNFGold(param *req.Param) bool {
	return strx.Equals(param.Message, "DNF", "dnf", "ゴールド", "金币", "比例")
}

func Translate(param *req.Param) bool {
	prefix := strx.SubStr(param.Message, -1, 2)
	return mapx.ContainsKey(handler.LangReflect, prefix)
}

func EroImagesSearch(param *req.Param) bool {
	return strx.StartWith(param.Message, "#")
}

func EroImages(param *req.Param) bool {
	return strx.Contains(param.Message, "色图", "涩图")
}

func SimpleReflect(param *req.Param) bool {
	return mapx.ContainsKey(BaseReflect, param.Message)
}

func Text(param *req.Param) bool {
	return strx.Equals(param.Message, "文章", "课文")
}

// SubscribePixiv 订阅p站色图
func SubscribePixiv(param *req.Param) bool {
	return strx.StartWith(param.Message, "订阅 ")
}

func UnSubscribePixiv(param *req.Param) bool {
	return strx.Equals(param.Message, "取消订阅")
}

func Health(param *req.Param) bool {
	return strx.Equals(param.Message, "ping")
}

func PasswdManage(param *req.Param) bool {
	return strx.StartWith(param.Message, "密码")
}

func ChatGPT(param *req.Param) bool {
	return true
}
