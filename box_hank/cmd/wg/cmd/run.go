package cmd

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/simonvetter/modbus"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/box"
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

	t1 := &box.ModbusTask{
		Client:  mc,
		Addr:    0x0,
		RegType: modbus.INPUT_REGISTER,
		UnitID:  0x21,
		Sender:  s,

		Code: "826_0-16",
		Type: common.WATER,

		Project: "1006",

		Ctx: context.Background(),
	}

	l := box.NewLoop(box.SeqJobs(t1))

	l.Run()

	return http.ListenAndServe(":10001", nil)
}
