package wecom

import (
	"fmt"

	"github.com/go-sphere/wecom-aibot-go-sdk/aibot"
	"github.com/twiglab/h2o/dbcli/ent"
)

type AIBot struct {
	Cli *ent.Client

	cmdMgr CommandManager

	wsClient *aibot.WSClient
}

func NewAIBot() *AIBot {
	bot := &AIBot{}

	bot.wsClient.OnMessageText(func(frame *aibot.WsFrame) {

		var msg aibot.TextMessage
		if err := aibot.ParseMessageBody(frame, &msg); err != nil {
			fmt.Println("解析消息失败:", err.Error())
			return
		}

		// fmt.Printf("收到文本: %s\n", msg.Text.Content)

		// cmd, _ := cm.Parser(topup.CmdCfg{Cli: cli}, strings.Split(msg.Text.Content, " "))
		// cmd.Do(context.Background())
		// client.Reply(frame, aibot.CreateMarkdownReplyBody(cmd.ToString()), "")
	})

	return bot
}
