package cmd

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/vigil"

	"github.com/spf13/cobra"
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

	rootLog()

	cli := entcli()
	defer func() { _ = cli.Close() }()

	hub := &vigil.Hub{
		DB: vigil.WithRecorder(dbx(cli)),
		// TSDB:   tdb(),
		Logger: serverLog(),
		WAL:    cwal(),
	}
	mcli := mqttcli2(func(mc mqtt.Client) {
		token := mc.Subscribe(GeneralDataTopic, 0x00, vigil.Handle(hub))
		if token.Wait() && token.Error() != nil {
			log.Fatal("mqtt - 订阅失败", token.Error())
		}
	})

	token := mcli.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatal("mqtt - connect error", token.Error())
	}

	go http.ListenAndServe(webaddr(), nil)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("quit...")
	return nil
}
