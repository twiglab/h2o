package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/twiglab/h2o/hank"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	s := &hank.Server{
		Addr: viper.GetString("server.addr"),
		Hub: &hank.Hub{
			DataLog: dataLog(),
			Sender:  sender(),
		},
		Logger: serverLog(),
	}

	go http.ListenAndServe(viper.GetString("web.addr"), nil)

	return s.Run()
}
