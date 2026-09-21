package cmd

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/chrgg"
	"github.com/twiglab/h2o/chrgg/orm"
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

	sl := serverLog()

	mcli := mqttcli()
	entc := entcli()

	act := chrgg.NewMQTTAction(mcli)

	svr := &chrgg.ChargeServer{
		Sender: act,
		DBx:    &orm.DBx{Cli: entc},
		Logger: sl,
	}
	t := mcli.Subscribe(common.GeneralDataTopic, 0x01, chrgg.HandleChange(svr))
	t.Wait()

	if err := t.Error(); err != nil {
		log.Fatal(err)
	}

	return http.ListenAndServe(webaddr(), nil)
}
