package cmd

import (
	"context"
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

// eyeCmd represents the eye command
var eyeCmd = &cobra.Command{
	Use:   "eye",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolEye()
	},
}

func init() {
	toolCmd.AddCommand(eyeCmd)
}

func toolEye() error {
	if deviceCode == "" {
		return errors.New("code is nil")
	}

	box, _ := boxData()
	topic := proto.MakeTopic(proto.DATA, box, deviceCode, "+")
	log.Println(topic)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mcli := mqttcli(
		nab.ClientID(box, randText()),
		topic,
		func(cli mqtt.Client, msg mqtt.Message) {
			fmt.Println(string(msg.Payload()))
		},
	)
	defer mcli.Disconnect(500)
	token := mcli.Connect()
	if err := mqtt.WaitTokenTimeout(token, 60*time.Second); err != nil {
		return err
	}

	<-ctx.Done()
	log.Println("quit...")
	return nil
}
