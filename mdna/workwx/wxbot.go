package workwx

import (
	"context"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/mdna"
	"github.com/xen0n/go-workwx/v2"
)

type WxBot struct {
	templ *template.Template
	whc   *workwx.WebhookClient
}

func NewWxBot(key string) WxBot {
	t := template.New("notice").Funcs(fm())

	template.Must(t.New("opt").Parse(optTempl))
	template.Must(t.New("chrgg").Parse(chrggTempl))

	return WxBot{
		whc:   workwx.NewWebhookClient(key),
		templ: t,
	}
}

func (w WxBot) StatusChangeNotice(ctx context.Context, n mdna.StatusChangeNoticeBody) error {
	var sb strings.Builder
	sb.Grow(2048)
	if err := w.templ.ExecuteTemplate(&sb, "opt", n); err != nil {
		return err
	}
	return w.whc.SendMarkdownV2Message(sb.String())
}
