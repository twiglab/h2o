package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/topup/wecom"
)

// wecombotCmd represents the wecombot command
var wecombotCmd = &cobra.Command{
	Use:   "wecombot",
	Short: "企业微信能耗智能机器人",
	Long:  `支持缴费，限额查询等功能`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return wecombot()
	},
}

func init() {
	rootCmd.AddCommand(wecombotCmd)
}

func wecombot() error {
	cli := entcli()
	defer cli.Close()

	wscli := wsclinet()
	defer wscli.Disconnect()

	cm := wecom.CmdMgr{
		Client:  cli,
		Context: context.Background(),

		Auth: fixGroup(),
	}
	cm.Init()

	wscli.OnMessageText(cm.TextMessageHandle(wscli))

	wscli.Connect()

	// 优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return nil
}
