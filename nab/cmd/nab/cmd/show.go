package cmd

import (
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
)

// showCmd represents the csv command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func toTable(e any) []string {
	switch v := e.(type) {
	case *ent.Dev:
		return orm.DevToStrings(v)
	case *ent.Cli:
		return orm.CliToStrings(v)
	}
	panic("no t")
}
