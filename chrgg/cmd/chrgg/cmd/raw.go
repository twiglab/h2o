/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/chrgg"
)

// rawCmd represents the raw command
var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return raw()
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)
}

func raw() error {

	_ = rootLog()

	c := mqttcli()

	c.SubscribeMultiple(topics(), chrgg.RawHandle())

	return http.ListenAndServe(webaddr(), nil)
}
