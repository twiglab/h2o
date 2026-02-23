/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"context"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/twiglab/h2o/hank"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
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
	rootCmd.AddCommand(csvCmd)
}

func csv() {

	load := viper.GetString("hank.abm.load")
	get := viper.GetString("hank.abm.get")
	list := viper.GetString("hank.abm.list")

	db := ddb()

	ctx := context.Background()

	if err := db.Load(ctx); err != nil {
		log.Fatal("load", err)
	}

	fmt.Println("load sql:", strings.TrimSpace(load))
	fmt.Println("----------------------------------")
	fmt.Println("get sql:", strings.TrimSpace(get))
	fmt.Println("----------------------------------")
	fmt.Println("list sql:", strings.TrimSpace(list))
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
	//sn,code,name,project,pos_code,building,floor_code,area_code,f1,f2,f3,f4,f5
	table.Header([]string{"sn", "code", "name", "project", "pos_code", "building", "floor_code", "area_code",
		"f1", "f2", "f3", "f4", "f5"})

	for _, emp := range rs {
		table.Append(emp)
	}

	table.Render()

}

func toTable(e any) []string {
	emp, ok := e.(hank.MetaData)
	if !ok {
		return []string{"Error: Invalid type"}
	}
	return emp.ToStrings()
}
