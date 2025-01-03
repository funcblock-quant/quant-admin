package config

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	AMap AMap // 这里配置对应配置文件的结构即可
	Lark LarkConfig
}

type AMap struct {
	Key string
}

type LarkConfig struct {
	Webhook string `yaml:"webhook"`
}
