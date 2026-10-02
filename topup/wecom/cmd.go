package wecom

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/twiglab/h2o/dbcli/ent"
)

type CmdCfg struct {
	Cli *ent.Client
}

type CommandMarker func(cfg CmdCfg, args ...string) (Commander, error)

type Commander interface {
	Do(context.Context) error
	ToString() string
}

type CommandManager struct {
	m map[string]CommandMarker
}

func NewCommandManager() CommandManager {
	m := make(map[string]CommandMarker)
	m["h"] = helpCmdFn
	m["?"] = helpCmdFn

	m["u"] = usageCmdFn
	m["p"] = posQueryCmdFn
	m["t"] = topCmdFn

	return CommandManager{m: m}
}

func (m CommandManager) Parser(cfg CmdCfg, input []string) (cmd Commander, err error) {
	args := slices.DeleteFunc(input, func(item string) bool {
		return strings.Contains(item, "@")
	})
	if cf, ok := m.m[args[0]]; ok {
		cmd, err = cf(cfg, args...)
		return
	}
	return hCmd, fmt.Errorf("not found cmd: %s", args[0])
}

type ErrorCmd struct {
	str string
}

func (e ErrorCmd) Error() string {
	return e.str
}
func (e ErrorCmd) Do(_ context.Context) error {
	return e
}

func (e ErrorCmd) ToString() string {
	return e.str
}
