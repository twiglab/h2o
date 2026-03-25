package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/chrgg"

	"context"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// feeCmd represents the csv command
var feeCmd = &cobra.Command{
	Use:   "fee",
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
	rootCmd.AddCommand(feeCmd)
}

func csv() {

	db, c := ddb()

	ctx := context.Background()

	if err := db.Load(ctx); err != nil {
		log.Fatal("load", err)
	}

	fmt.Println("load sql:", strings.TrimSpace(c.LoadSQL))
	fmt.Println("----------------------------------")
	fmt.Println("get sql:", strings.TrimSpace(c.GetSQL))
	fmt.Println("----------------------------------")
	fmt.Println("list sql:", strings.TrimSpace(c.ListSQL))
	fmt.Println("----------------------------------")

	rs, err := db.List(ctx)
	if err != nil {
		log.Fatal("list ", err)
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
	table.Header([]string{"code", "fee_fen", "pos_code"})

	for _, emp := range rs {
		table.Append(emp)
	}

	table.Render()
}

func toTable(e any) []string {
	emp, ok := e.(chrgg.AloneRuler)
	if !ok {
		return []string{"Error: Invalid type"}
	}
	return emp.ToStrings()
}
