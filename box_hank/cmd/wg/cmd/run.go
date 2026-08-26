package cmd

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"time"

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

	// 2025100021	F5-02
	t1 := &box.ModbusTask{
		Client:  mc,
		Addr:    0x0,
		RegType: modbus.INPUT_REGISTER,
		UnitID:  0x21,
		Sender:  s,

		Code: "826_0-21",
		Type: common.WATER,

		Project: "1006",

		Ctx: context.Background(),
	}

	// 2024082045	F5-01
	t2 := &box.ModbusTask{
		Client:  mc,
		Addr:    0x0,
		RegType: modbus.INPUT_REGISTER,
		UnitID:  0x45,
		Sender:  s,

		Code: "826_0-45",
		Type: common.WATER,

		Project: "1006",

		Ctx: context.Background(),
	}

	// 2026050057	F5-02
	t3 := &box.ModbusTask{
		Client:  mc,
		Addr:    0x0,
		RegType: modbus.INPUT_REGISTER,
		UnitID:  0x57,
		Sender:  s,

		Code: "826_0-57",
		Type: common.WATER,

		Project: "1006",

		Ctx: context.Background(),
	}

	l := box.NewLoop(30*time.Second, box.SeqJobs(3*time.Second, t1, t2, t3))

	l.Run()

	return http.ListenAndServe(":10001", nil)
}
