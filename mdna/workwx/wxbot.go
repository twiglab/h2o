package workwx

import (
	"context"
	"strings"
	"text/template"
	"time"

	"github.com/twiglab/h2o/mdna"
	"github.com/xen0n/go-workwx/v2"
)

type WxBot struct {
	BotKey string
	Tmpl   *template.Template
}

func NewWxBot(key string) WxBot {
	tmpl := template.New("notice").Funcs(fm())
	tmpl.New("opt").Parse(optTempl)
	tmpl.New("chrgg").Parse(chrggTempl)

	return WxBot{
		BotKey: key,
		Tmpl:   tmpl,
	}
}

func (w WxBot) Notice(ctx context.Context, n mdna.Notice) error {
	wc := workwx.NewWebhookClient(w.BotKey)

	now := time.Now()

	var sb strings.Builder
	sb.Grow(2048)
	if err := w.Tmpl.ExecuteTemplate(&sb, n.Type, NotictBox{
		Now:    now,
		Notice: n,
	}); err != nil {
		return err
	}

	return wc.SendMarkdownV2Message(sb.String())
}
