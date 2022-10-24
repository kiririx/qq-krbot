package env

import "github.com/kiririx/krutils/confx"

var Conf map[string]string

func init() {
	conf, err := confx.ResolveProperties("./config.properties")
	if err != nil {
		panic("配置文件读取失败")
	}
	Conf = conf
}
