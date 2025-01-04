package jobs

import (
	"fmt"
	"quanta-admin/notification/lark"
	"time"
)

// InitJob
// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func InitJob() {
	jobList = map[string]JobExec{
		"ExamplesOne":                 ExamplesOne{},
		"MonitorArbitrageOpportunity": MonitorArbitrageOpportunity{}, //监控套利机会定时任务
		// ...
	}
}

// ExamplesOne
// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct {
}

func (t ExamplesOne) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore ExamplesOne exec success"
	// TODO: 这里需要注意 Examples 传入参数是 string 所以 arg.(string)；请根据对应的类型进行转化；
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			//testLarkNotification()
			fmt.Println("string", arg.(string))
			fmt.Println(str, arg.(string))
		} else {
			fmt.Println("arg is nil")
			fmt.Println(str, "arg is nil")
		}
		break
	}

	return nil
}

func testLarkNotification() {
	fmt.Println("测试lark通知")
	notification := lark.NewLarkTextNotification(new(string), "测试")
	notification.SendNotification()
}

// MonitorArbitrageOpportunity
type MonitorArbitrageOpportunity struct {
}

func (t MonitorArbitrageOpportunity) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore MonitorArbitrageOpportunity exec success"

	fmt.Println(str)
	return nil
}
