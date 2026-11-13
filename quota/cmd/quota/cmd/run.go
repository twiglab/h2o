package cmd

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/quota"
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

	entc := entcli()

	defer func() { _ = entc.Close() }()

	db := &quota.DBx{Cli: entc}

	wl := whiteList()

	mcli := mqttcli(func(cli mqtt.Client) {
		log.Println("mqtt - onConnect")

		svr := &quota.ChargeServer{
			DBx: db,

			MCli:   cli,
			Sender: quota.NewMQTTAction(cli),

			WhiteList: wl,

			WAL:    cwal(),
			Logger: sl,
		}

		token := cli.Subscribe(quota.GeneralDataTopic, 0x00, quota.MsgHandle(svr))
		if token.Wait() && token.Error() != nil {
			log.Fatal("mqtt - 订阅失败", token.Error())
		}

	})

	token := mcli.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatal("mqtt - connect error", token.Error())
	}

	return http.ListenAndServe(webaddr(), nil)
}
