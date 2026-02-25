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
	ep := hank.NewElectricityPacket()

	s := &hank.Server{
		Addr: viper.GetString("hank.server.addr"),
		Hub: &hank.Hub{
			WAL:    walLog(),
			Sender: sender(),
			EP:     ep,
		},
		Logger: serverLog(),
		Enh:    enh(),
	}

	http.HandleFunc("/eyes/all", hank.EyesAll(ep))
	go http.ListenAndServe(viper.GetString("hank.web.addr"), nil)

	return s.Run()
}
