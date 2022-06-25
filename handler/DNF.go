package handler

import (
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/json_util"
	"time"
)

type DNF struct {
}

type DNFGold struct {
	Scale     float64
	TradeType string
}

func DNFHandler() *DNF {
	return &DNF{}
}

func (d *DNF) Gold() ([]DNFGold, error) {
	url := "http://www.uu898.com/ashx/GameRetail.ashx?act=a001&g=95&a=2335&s=25080&c=-3&cmp=-1&_t=1639304944162"
	resp, err := http_util.Client(time.Second*4).GetJSON(url, nil)
	if err != nil {
		return nil, err
	}
	data, err := json_util.JSON2Map(resp)
	if err != nil {
		return nil, err
	}
	gold := data["list"].(map[string]interface{})["datas"].([]interface{})
	goldList := make([]DNFGold, 0)
	for _, v := range gold {
		goldList = append(goldList, DNFGold{
			Scale:     v.(map[string]interface{})["Scale"].(float64),
			TradeType: v.(map[string]interface{})["TradeType"].(string),
		})
	}
	return goldList, nil
}
