package cmd

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/simonvetter/modbus"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/box"
	"github.com/twiglab/h2o/box/ocg"
	"github.com/twiglab/h2o/pkg/common"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
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

	mc := modbusClient()
	defer mc.Close()

	s := sender()

	t1 := &ocg.TaskX{
		Client:  mc,
		Addr:    2,
		RegType: modbus.INPUT_REGISTER,
		Sender:  s,

		Code: "PT-1-IN",
		Type: common.ELECTRICITY,

		Project: "1006",
	}

	t2 := &ocg.TaskX{
		Client:  mc,
		Addr:    16,
		RegType: modbus.INPUT_REGISTER,
		Sender:  s,

		Code: "PT-2-IN",
		Type: common.ELECTRICITY,

		Project: "1006",
	}

	cron := box.NewCron()

	cron.AddFunc("@every 15m", func() {
		var err error
		if err = ocg.TaskChain(context.Background(), t1, t2); err != nil {
			slog.Error("task error", slog.Any("error", err))
		}
	})
	cron.Start()

	return http.ListenAndServe(":10000", nil)
}
