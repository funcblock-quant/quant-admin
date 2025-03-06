package lark

import (
	"fmt"
	ext "quanta-admin/config"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-lark/lark"
)

type LarkRobotAlert struct {
	webhook string
	secret  string
}

func NewLarkRobotAlert(larkBotConf ext.Extend) *LarkRobotAlert {
	return &LarkRobotAlert{
		webhook: larkBotConf.Lark.Webhook,
		secret:  larkBotConf.Lark.Secret,
	}
}

func (a LarkRobotAlert) SendLarkAlert(text string) error {
	bot := lark.NewNotificationBot(a.webhook)
	secret := a.secret
	message := lark.NewMsgBuffer(lark.MsgText)
	if secret != "" {
		message = message.WithSign(secret, time.Now().Unix())
	}

	resp, err := bot.PostNotificationV2(message.Text(text).Build())
	if err != nil {
		log.Errorf("bot.PostNotificationV2 error: %v", err)
		return err
	}
	//log.Infof("resp: %+v", resp.Code)
	if resp.Code != 0 {
		return fmt.Errorf("resp.Code: %v, error message:%s", resp.Code, resp.Msg)
	}
	return nil
}
