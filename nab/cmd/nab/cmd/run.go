package cmd

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"

	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/orm/ent"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行采集服务",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run() error {

	_ = rootLog()
	logger := serverLog()
	mqttLog()
	onofflog := onoffLog()

	p := pool()
	p.Release()

	boxCode, project := boxData()

	box := &nab.Box{
		Code:    boxCode,
		Project: project,
		DB:      db(boxCode),
	}

	if err := box.Load(); err != nil {
		return err
	}
	defer box.Disconnect()

	agent := &nab.Agent{
		Box:         box,
		Logger:      logger,
		OnoffLogger: onofflog,
		Pool:        p,
	}

	mcli := mqttcli(
		nab.ClientID(boxCode, nab.NAB),
		nab.SubscriptTopic(boxCode),
		nab.AgentHandle(agent),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lps := loops(box, func(data []*ent.Dev) nab.JobV2 {
		return nab.CollectTask{
			Box:      box,
			Data:     data,
			Logger:   logger,
			Pipeline: nab.NewMQTTAction(mcli),
			Delay:    taskInterval(),
		}
	})

	token := mcli.Connect()
	if err := mqtt.WaitTokenTimeout(token, 60*time.Second); err != nil {
		slog.Error("mqtt connect error", slog.Any("error", err))
		return err
	}
	defer mcli.Disconnect(10000) //10s

	lspCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	lps.Run(lspCtx)

	<-ctx.Done()

	log.Println("quit...")
	return nil
}
