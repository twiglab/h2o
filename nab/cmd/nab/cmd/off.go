package cmd

import (
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/proto"
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolOnoff(proto.OFF)
	},
}

func init() {
	toolCmd.AddCommand(offCmd)
}
