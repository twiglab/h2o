package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
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

	mcli := mcli2(vigil.RawHandle())

	token := mcli.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatal("mqtt - connect error", token.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("quit...")
	return nil
}

const GeneralDataTopic = "h2o/data/#"

const f = "20060102150405"

func Ts(now time.Time) string {
	return now.Format(f)
}

func mcli2(h mqtt.MessageHandler) mqtt.Client {
	broker := viper.GetString("h2o.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	ts := Ts(time.Now())
	clientID := "eye@" + ts

	opts := mqtt.NewClientOptions()
	opts.SetClientID(clientID)
	opts.AddBroker(broker)
	opts.SetOrderMatters(false)
	opts.SetAutoReconnect(true)

	opts.SetOnConnectHandler(func(mc mqtt.Client) {
		token := mc.Subscribe(topic, 0x00, h)
		if token.Wait() && token.Error() != nil {
			slog.Error("mqtt - subscribe error", slog.Any("error", token.Error()))
		}
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		slog.Error("mqtt - connect lost", slog.Any("error", err))
	})

	client := mqtt.NewClient(opts)
	log.Printf("mqtt - broker: %s, clientID: %s\n", broker, vigil.APP_ID)
	return client
}
