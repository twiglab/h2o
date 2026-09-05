/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/pkg/common"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eyeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eyeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func eye(cmd *cobra.Command, args []string) error {

	rLog()

	mcli := mcli()
	token := mcli.Subscribe(common.ElectricityTopic, 0x0, vigil.RawHandle())
	token.Wait()

	return http.ListenAndServe(":10020", nil)
}

func mcli() mqtt.Client {
	broker := viper.GetString("vigil.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	cli, err := vigil.NewMQTTClient("eye", broker)
	if err != nil {
		log.Fatal(fmt.Errorf("mc err: %w", err))
	}
	return cli
}

func rLog() *slog.Logger {
	log := clog.NewLog("", slog.LevelDebug)
	slog.SetDefault(log)
	return log
}
