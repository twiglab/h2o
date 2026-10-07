package wecom

import (
	"fmt"

	"github.com/go-sphere/wecom-aibot-go-sdk/aibot"
)

func NewWsClient(botID, botSecret string) *aibot.WSClient {

	/*
		if botID == "" || botSecret == "" {
			fmt.Println("请设置环境变量 WECOM_BOT_ID 和 WECOM_BOT_SECRET")
			return
		}
	*/

	client := aibot.NewWSClient(aibot.WSClientOptions{
		BotID:  botID,
		Secret: botSecret,
	})

	// 设置事件处理
	client.OnConnected(func() {
		fmt.Println("连接已建立")
	})

	client.OnAuthenticated(func() {
		fmt.Println("认证成功")
	})

	client.OnDisconnected(func(reason string) {
		fmt.Println("连接断开:", reason)
	})

	// 监听进入会话事件（发送欢迎语）
	client.OnEventEnterChat(func(frame *aibot.WsFrame) {
		welcomeBody := aibot.CreateTextReplyBody("您好！我是智能助手，有什么可以帮您的吗？")
		client.ReplyWelcome(frame, welcomeBody)
	})

	return client
}
