package resp

import (
	"qq-krbot/req"
	"testing"
)

func TestChatGPT(t *testing.T) {
	for {
		ChatGPT(&req.Param{
			PostType:    "",
			Message:     "",
			MessageType: "",
			SubType:     "",
			GroupId:     0,
			RawMessage:  "",
			UserId:      0,
			CQ:          "",
		})
	}

}
