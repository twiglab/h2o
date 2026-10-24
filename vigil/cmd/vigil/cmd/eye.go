package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/vigil"
)

// eyeCmd represents the eye command
var eyeCmd = &cobra.Command{
	Use:   "eye",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: eye,
}

func init() {
	rootCmd.AddCommand(eyeCmd)
	eyeCmd.Flags().StringVarP(&topic, "topic", "t", GeneralDataTopic, "Help message for toggle")
}

var topic string

func eye(cmd *cobra.Command, args []string) error {
	fmt.Println("topic ", topic)

	clog.ConsoleLog(slog.LevelDebug)

	mcli := mcli()
	token := mcli.Subscribe(topic, 0x0, vigil.RawHandle())
	token.Wait()

	return http.ListenAndServe(":10020", nil)
}

func mcli() mqtt.Client {
	broker := viper.GetString("vigil.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	ts := Ts(time.Now())
	cli, err := vigil.NewMQTTClient("eye@"+ts, broker)
	if err != nil {
		log.Fatal(fmt.Errorf("mc err: %w", err))
	}
	return cli
}
