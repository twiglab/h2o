package cmd

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/proto"
)

// toolCmd represents the tool command
var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	rootCmd.AddCommand(toolCmd)
	toolCmd.PersistentFlags().StringVar(&deviceCode, "code", "+", "device code")
}

var deviceCode string

func toolOnoff(onoff string) error {
	if deviceCode == "" || deviceCode == "+" || deviceCode == "#" {
		return errors.New("code is nil")
	}

	box, _ := boxData()
	topic := proto.MakeTopic(proto.ONOFF, box, deviceCode, "+")
	log.Println(topic)
	log.Println(rand.Text())
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mcli := mqttcli(
		nab.ClientID(box, randText()),
		topic,
		func(cli mqtt.Client, msg mqtt.Message) {
			fmt.Println(msg.Topic())
		},
	)
	defer mcli.Disconnect(500)
	token := mcli.Connect()
	if !token.WaitTimeout(60 * time.Second) {
		return errors.New("mqtt connect timeout")
	}
	if err := token.Error(); err != nil {
		return err
	}

	act := nab.NewMQTTAction(mcli)

	err := act.SendData(ctx, nab.OnoffLite{
		Box:  box,
		Code: deviceCode,
		Op:   onoff,
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	log.Println("quit...")
	return nil
}
