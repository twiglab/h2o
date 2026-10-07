package wecom

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/go-sphere/wecom-aibot-go-sdk/aibot"
	"github.com/twiglab/h2o/dbcli/ent"
)

type Global struct {
	Client   *ent.Client
	Template *template.Template
	Context  context.Context

	Frame       *aibot.WsFrame
	TextMessage aibot.TextMessage

	Auth *FixUserGroup
}

type CmdMakeFn func(cfg Global, args ...string) (Commander, error)

type Commander interface {
	Do(context.Context) error
	ToString() string
}

type CmdMgr struct {
	Client  *ent.Client
	Auth    *FixUserGroup
	Context context.Context

	templ *template.Template
	m     map[string]CmdMakeFn
}

func (c *CmdMgr) Init() {
	c.templ = BuildTemplate()

	m := make(map[string]CmdMakeFn)
	m["h"] = helpCmdFn
	m["?"] = helpCmdFn

	m["u"] = usageCmdFn
	m["p"] = posQueryCmdFn
	m["t"] = topCmdFn

	m["c"] = chargeCmdFn

	c.m = m
}

func (c CmdMgr) TextMessageHandle(wscli *aibot.WSClient) func(*aibot.WsFrame) {
	return func(frame *aibot.WsFrame) {
		var msg aibot.TextMessage
		if err := aibot.ParseMessageBody(frame, &msg); err != nil {
			fmt.Println("解析消息失败:", err.Error())
			return
		}

		fmt.Printf("收到文本: %s\n", msg.Text.Content)

		cfg := Global{
			Client:   c.Client,
			Template: c.templ,
			Context:  c.Context,

			Frame:       frame,
			TextMessage: msg,

			Auth: c.Auth,
		}

		args := strings.FieldsFunc(msg.Text.Content, isField)

		cmd, err := c.Parser(cfg, args)
		if err != nil {
			wscli.Reply(frame, aibot.CreateMarkdownReplyBody(err.Error()), "")
			return
		}

		if err := cmd.Do(context.Background()); err != nil {
			wscli.Reply(frame, aibot.CreateMarkdownReplyBody(err.Error()), "")
			return
		}
		wscli.Reply(frame, aibot.CreateMarkdownReplyBody(cmd.ToString()), "")
	}
}

func (c CmdMgr) Parser(cfg Global, input []string) (cmd Commander, err error) {
	args := slices.DeleteFunc(input, func(item string) bool {
		return strings.Contains(item, "@")
	})

	cc := strings.ToLower(args[0])
	if cf, ok := c.m[cc]; ok {
		cmd, err = cf(cfg, args...)
		return
	}
	return helpCmdFn(cfg, input...)
}

func isField(r rune) bool {
	switch r {
	case ' ':
		return true
	case '\t':
		return true
	case ',':
		return true
	}
	return false
}
