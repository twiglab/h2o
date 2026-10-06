package wecom

import (
	"context"
	"strings"
	"text/template"
)

type helpCmd struct {
	Template *template.Template
}

func (h helpCmd) Do(ctx context.Context) error {
	return nil
}

func (h helpCmd) ToString() string {
	var sb strings.Builder
	_ = h.Template.ExecuteTemplate(&sb, "help", nil)
	return sb.String()
}

func helpCmdFn(cfg Global, args ...string) (Commander, error) {
	help := helpCmd{Template: cfg.Template}
	return help, nil
}
