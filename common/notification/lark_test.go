package lark

import (
	"fmt"
	ext "quanta-admin/config"
	"testing"
	"time"
)

func TestLarkNotification(t *testing.T) {
	larkRobotConf := ext.Extend{
		Lark: ext.LarkConfig{
			Webhook: "https://open.larksuite.com/open-apis/bot/v2/hook/2264a2cd-0acf-4a11-85a0-5fab1bf510be",
			Secret:  "PC8UjZAbQw3IA4za7bIKOg",
		},
	}
	larkRobot := NewLarkRobotAlert(larkRobotConf)

	// 1. 预警类消息
	warningMessage := fmt.Sprintf(`
⚠️ 风控预警通知
	策略ID: %s
	实例ID: %s
	触发条件: %s
	当前值: %s %s
	通知时间: %s

🔔 该预警不会影响交易，仅供参考。请关注交易风险。
	`, "strategy123", "instance456", "单笔亏损超过 10U", "12", "USDT", time.Now().Format(time.RFC3339))

	// 2. 暂停单币交易消息
	pauseSingleMessage := fmt.Sprintf(`
🚨 风控触发：暂停 %s 交易
	策略ID: %s
	实例ID: %s
	触发条件: %s
	当前值: %s %s
	暂停时间: %d 秒
	恢复方式: %s
	通知时间: %s

❗ 请立即检查策略，并决定是否手动恢复交易。
	`, "USDT", "strategy123", "instance456", "单笔亏损超过 20U", "22", "USDT", 3600, "手动恢复", time.Now().Format(time.RFC3339))

	// 3. 全局暂停交易消息
	globalPauseMessage := fmt.Sprintf(`
🛑 交易系统已全局暂停
	触发策略: %s
	实例ID: %s
	触发条件: %s
	当前值: %s
	暂停原因: %s
	恢复方式: %s
	通知时间: %s

🚨 全局交易已暂停，请立即检查风险，并决定恢复方案。
	`, "strategy123", "instance456", "单笔亏损超过 50U", "55", "单笔亏损超限", "手动恢复", time.Now().Format(time.RFC3339))

	// 调用 Lark 机器人发送预警消息
	err := larkRobot.SendLarkAlert(warningMessage)
	if err != nil {
		t.Errorf("SendLarkAlert failed for warning, err:%v", err)
	}

	// 调用 Lark 机器人发送暂停单币交易消息
	err = larkRobot.SendLarkAlert(pauseSingleMessage)
	if err != nil {
		t.Errorf("SendLarkAlert failed for pause single, err:%v", err)
	}

	// 调用 Lark 机器人发送全局暂停交易消息
	err = larkRobot.SendLarkAlert(globalPauseMessage)
	if err != nil {
		t.Errorf("SendLarkAlert failed for global pause, err:%v", err)
	}

}
