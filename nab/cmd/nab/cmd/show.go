package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/nab/orm/idb"

	"context"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
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
	Run: func(cmd *cobra.Command, args []string) {
		csv()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func csv() {

	cli, err := orm.NewIDB("idb/dev.csv", "idb/cli.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	rs, err := cli.Dev.Query().All(ctx)
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithStringer(toTable),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.On},
				Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}},
		}),
	)

	table.Header(idb.DEVICE_TABLE_COLUMNS)

	for _, r := range rs {
		fmt.Println(r)
		table.Bulk(r)
	}

	table.Render()
}

func toTable(e any) []string {
	emp, ok := e.(*ent.Dev)
	if !ok {
		return []string{"Error: Invalid type"}
	}
	return orm.DevToStrings(emp)
}
