package main

import (
	"regexp"
	"strings"
)

type Param struct {
	PostType    string `json:"post_type"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	GroupId     int64  `json:"group_id"`
	RawMessage  string `json:"raw_message"`
	UserId      int64  `json:"user_id"`
	CQ          string
}

// [CQ:at,qq=1105624883] 时尚
func (p *Param) Parse() {
	if p.PostType == "message" {
		regex := `\[CQ:at,qq=(\d+)\]`
		reg, err := regexp.Compile(regex)
		if err != nil {
			return
		}
		matchedStr := reg.FindString(p.Message)
		cqPrefix := strings.Split(matchedStr, ",")[0]
		CQ := strings.Replace(cqPrefix, "[CQ:", "", -1)
		if strings.HasPrefix(p.Message, "[CQ:at,qq=") {
			p.CQ = CQ
			p.Message = strings.TrimPrefix(p.Message, matchedStr+" ")
		}
	}
}
