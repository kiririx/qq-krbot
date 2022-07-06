package cond

import (
	"github.com/kiririx/krutils/map_util"
	"github.com/kiririx/krutils/str_util"
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
	return str_util.Contains(param.Message, "ヘルプ", "帮助", "help", "?", "？")
}

func DNFGold(param *req.Param) bool {
	return str_util.Contains(param.Message, "DNF", "dnf", "ゴールド", "金币", "比例")
}

func Translate(param *req.Param) bool {
	prefix := str_util.SubStr(param.Message, -1, 2)
	return map_util.ContainsKey(handler.LangReflect, prefix)
}

func EroImagesSearch(param *req.Param) bool {
	return str_util.StartWith(param.Message, "#")
}

func EroImages(param *req.Param) bool {
	return str_util.Contains(param.Message, "色图", "涩图")
}

func SimpleReflect(param *req.Param) bool {
	return map_util.ContainsKey(BaseReflect, param.Message)
}

func Text(param *req.Param) bool {
	return str_util.Equals(param.Message, "文章", "课文")
}

// SubscribePixiv 订阅p站色图
func SubscribePixiv(param *req.Param) bool {
	return str_util.StartWith(param.Message, "订阅 ")
}
