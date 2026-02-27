/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"context"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"code", "fee_fen", "pos_code"})

	for _, emp := range rs {
		table.Append(emp.ToStrings())
	}

	table.Render()
}
