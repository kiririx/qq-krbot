package trigger

import (
	"qq-krbot/cond"
	"qq-krbot/req"
	"qq-krbot/resp"
)

var (
	Triggers []Trigger
)

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
	addTrigger("at", cond.Help, resp.Help)
	addTrigger("at", cond.DNFGold, resp.DNFGold)
	addTrigger("at", cond.Translate, resp.Translate)
	addTrigger("at", cond.EroImagesSearch, resp.EroImagesSearch)
	addTrigger("at", cond.EroImages, resp.EroImages)
	addTrigger("at", cond.Text, resp.Text)
	addTrigger("pr", cond.SubscribePixiv, resp.SubscribePixiv)
	addTrigger("pr", cond.UnSubscribePixiv, resp.UnSubscribePixiv)
	addTrigger("pr", cond.Health, resp.Health)
}
