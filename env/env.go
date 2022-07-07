package env

import "github.com/kiririx/krutils/conf_util"

var Conf map[string]string

func init() {
	conf, err := conf_util.ResolveProperties("./config.properties")
	if err != nil {
		panic("配置文件读取失败")
	}
	Conf = conf
}
