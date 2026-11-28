package cmd

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/nab/orm/idb"
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showCli()
	},
}

func init() {
	showCmd.AddCommand(cliCmd)
}

func showCli() error {
	box, _ := boxData()

	db := db(box)
	defer db.Close()

	ctx := context.Background()

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.On},
				Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}},
		}),
	)

	rs, err := db.AllCli(ctx)
	if err != nil {
		return err
	}

	table.Header(idb.CLIENT_TABLE_COLUMNS)

	for _, r := range rs {
		table.Append(toTable(r))
	}

	return table.Render()
}
