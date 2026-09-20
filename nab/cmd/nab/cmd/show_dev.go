/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/nab/orm/idb"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showDev()
	},
}

func init() {
	showCmd.AddCommand(devCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func showDev() error {
	cli := db()
	defer cli.Close()

	ctx := context.Background()

	table := tablewriter.NewTable(os.Stdout,
		// tablewriter.WithStringer(toTable),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.On},
				Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row:   tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}},
			Debug: true,
		}),
	)

	rs, err := cli.AllDev(ctx)
	if err != nil {
		log.Fatal(err)
	}

	table.Header(idb.DEVICE_TABLE_COLUMNS)

	for _, r := range rs {
		table.Append(toTable(r))
	}

	return table.Render()
}
