package main

import (
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/map_util"
	"github.com/kiririx/krutils/str_util"
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
)

var (
	LangCorr map[int]LangDict
)

func init() {
	LangCorr = make(map[int]LangDict)
	LangCorr[cj] = initLangDict("zh-CHS", "ja")
	LangCorr[ce] = initLangDict("zh-CHS", "en")
	LangCorr[jc] = initLangDict("ja", "zh-CHS")
	LangCorr[je] = initLangDict("ja", "en")
	LangCorr[ec] = initLangDict("en", "zh-CHS")
	LangCorr[ej] = initLangDict("en", "ja")
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
	if !map_util.ContainsKey(LangCorr, tranType) {
		return ""
	}
	langDict := LangCorr[tranType]
	var input string
	_len := str_util.Len(text)
	if _len > 20 {
		input = str_util.SubStr(text, -1, 10) + str_util.ToStr(_len) + str_util.SubStr(text, _len-10, -1)
	} else {
		input = text
	}
	salt := str_util.ToStr(time.Now().UnixMilli())
	curtime := str_util.ToStr(time.Now().Unix())
	params := map[string]string{
		"q":        text,
		"from":     langDict.From,
		"to":       langDict.To,
		"appKey":   AppKey,
		"salt":     salt,
		"sign":     algo_util.Sha256(AppKey + input + salt + curtime + AppSecret),
		"signType": "v3",
		"curtime":  curtime,
	}
	resp, err := http_util.Client(time.Second*4).PostFormGetJSON(YDAPIURL, params)
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
