package lark

import (
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-lark/lark"
	ext "quanta-admin/config"
)

// LarkNotification lark通知消息体
type Notification struct {
	webHookUrl string
	msgText    string
	msg        string
	title      string
}

func NewLarkTextNotification(webHookUrl *string, msg string) *Notification {
	notification := Notification{
		msg:     msg,
		msgText: lark.MsgText,
	}
	if webHookUrl != nil && *webHookUrl != "" {
		notification.webHookUrl = *webHookUrl
		//log.Infof("lark webhook url: %s", *webHookUrl)
	} else {
		//log.Infof("default lark webhook url: %s", ext.ExtConfig.Lark.Webhook)
		notification.webHookUrl = ext.ExtConfig.Lark.Webhook
	}
	return &notification
}

func (n *Notification) SendNotification() error {
	var err error
	bot := lark.NewNotificationBot(n.webHookUrl)
	_, err = bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(n.msg).Build())
	if err != nil {
		log.Error("lark bot post notification failed", err)
		return err
	}
	log.Info("lark bot post notification success")
	return nil
}

func (n *Notification) SendCardNotification() error {
	// TODO 发送卡片消息
	return nil
}
