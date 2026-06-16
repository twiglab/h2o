package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// csvCmd represents the csv command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var VersionOverride = ""

func version() {
	if VersionOverride != "" {
		fmt.Println(VersionOverride)
		return
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			fmt.Println(info.Main.Version)
			return
		}
	}
	fmt.Println("(unknown)")
}
