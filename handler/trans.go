package handler

import (
	"github.com/kiririx/krutils/algox"
	"github.com/kiririx/krutils/httpx"
	"github.com/kiririx/krutils/mapx"
	"github.com/kiririx/krutils/strx"
	"log"
	"time"
)

const (
	cj = iota
	ce
	jc
	je
	ec
	ej
	// YDAPIURL 有道词典API
	YDAPIURL = "https://openapi.youdao.com/api"
	// AppKey 有道APP_KEY
	AppKey = "6a291bcc994b9a8b"
	// AppSecret 有道APP_SECRET
	AppSecret = "XUVMpL79kRKgPRkFAcDGCPDKjihw9aO1"
)

var (
	LangCorr    map[int]LangDict
	LangReflect map[string]int
)

func init() {
	LangCorr = make(map[int]LangDict)
	LangCorr[cj] = initLangDict("zh-CHS", "ja")
	LangCorr[ce] = initLangDict("zh-CHS", "en")
	LangCorr[jc] = initLangDict("ja", "zh-CHS")
	LangCorr[je] = initLangDict("ja", "en")
	LangCorr[ec] = initLangDict("en", "zh-CHS")
	LangCorr[ej] = initLangDict("en", "ja")

	// 语言翻译映射
	LangReflect = make(map[string]int)
	LangReflect["中日"] = cj
	LangReflect["中英"] = ce
	LangReflect["英中"] = ec
	LangReflect["日中"] = jc
	LangReflect["日英"] = je
	LangReflect["英日"] = ej
}

type LangDict struct {
	From string
	To   string
}

func initLangDict(from, to string) LangDict {
	return LangDict{
		From: from,
		To:   to,
	}
}

func Translate(tranType int, text string) string {
	if !mapx.ContainsKey(LangCorr, tranType) {
		return ""
	}
	langDict := LangCorr[tranType]
	var input string
	_len := strx.Len(text)
	if _len > 20 {
		input = strx.SubStr(text, -1, 10) + strx.ToStr(_len) + strx.SubStr(text, _len-10, -1)
	} else {
		input = text
	}
	salt := strx.ToStr(time.Now().UnixMilli())
	curtime := strx.ToStr(time.Now().Unix())
	params := map[string]string{
		"q":        text,
		"from":     langDict.From,
		"to":       langDict.To,
		"appKey":   AppKey,
		"salt":     salt,
		"sign":     algox.Sha256(AppKey + input + salt + curtime + AppSecret),
		"signType": "v3",
		"curtime":  curtime,
	}
	resp, err := httpx.Client().PostFormGetJSON(YDAPIURL, params)
	log.Println("params => ", params)
	log.Println(resp, err)
	if err != nil || resp["errorCode"] != "0" {
		log.Println(err)
		return ""
	}
	result := resp["translation"].([]any)
	var resultStr string
	for _, v := range result {
		resultStr = resultStr + v.(string) + ";"
	}
	return resultStr
}
